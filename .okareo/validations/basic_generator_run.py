from sklearn.metrics import classification_report
from okareo import Okareo
from okareo_api_client.models import ScenarioSetCreate, SeedData
from utils import read_in_data_points, run_pred, get_tfidf_featurizer, get_rf_classifier, download_scenario_data_points, parse_sceenario_json_data, check_basic_assertion, get_svm_classifier, get_sbert_featurizer
import random
import string
import os

OKAREO_API_KEY = os.environ["OKAREO_API_KEY"]
okareo = Okareo(OKAREO_API_KEY)
DELTA_DIFF = 0.02  # Large delta is more strict and better
file_path = os.path.join(
    os.path.dirname(__file__), "datasets/sentiment/binary_small.txt"
)


def run_basic_generator(run_name, scenario_type):
    ########################## Read In Data Points ####################################
    data = read_in_data_points(file_path)
    seed_data = [SeedData(input_=x, result=y) for x, y in data]
    print(f'sample size: {len(seed_data)}')

    ########################## Create/Download Scenario Set ###########################
    random_string = ''.join(random.choices(string.ascii_letters, k=10))
    max_num_samples = 5

    scenario_set_create = ScenarioSetCreate(
        name=f"GitHub {run_name} Scenario Categories Baseline " +
        random_string,
        number_examples=max_num_samples,
        generation_type=scenario_type,
        seed_data=seed_data
    )

    scenario = okareo.create_scenario_set(scenario_set_create)
    scenario_id = scenario.app_link.split('/')[-1]

    response = download_scenario_data_points(scenario_id, OKAREO_API_KEY)
    expanded_data_X, expanded_data_Y = parse_sceenario_json_data(response.text)

    print(f'generated sample size: {len(expanded_data_X)}')

    ########################## Running Some Basic Models ##############################
    train_X, train_Y = zip(*data)
    seed_test_X, seed_test_Y = list(train_X), list(train_Y)
    expanded_test_X, expanded_test_Y = expanded_data_X, expanded_data_Y

    ### Loading Featurizers ###
    tfidf_featurizer = get_tfidf_featurizer(train_X)
    sbert_featurizer = get_sbert_featurizer()

    ### Creating Models ###
    tfidf_rf_classifier = get_rf_classifier(
        tfidf_featurizer.transform(train_X), train_Y)
    tfidf_svm_classifier = get_svm_classifier(
        tfidf_featurizer.transform(train_X), train_Y)
    sbert_rf_classifier = get_rf_classifier(
        sbert_featurizer.encode(train_X), train_Y)
    sbert_svm_classifier = get_svm_classifier(
        sbert_featurizer.encode(train_X), train_Y)

    ### Running models ###
    tfidf_rf_seed_pred = run_pred(
        tfidf_rf_classifier, tfidf_featurizer.transform(seed_test_X))
    tfidf_rf_expanded_pred = run_pred(
        tfidf_rf_classifier, tfidf_featurizer.transform(expanded_test_X))

    tfidf_svm_seed_pred = run_pred(
        tfidf_svm_classifier, tfidf_featurizer.transform(seed_test_X))
    tfidf_svm_expanded_pred = run_pred(
        tfidf_svm_classifier, tfidf_featurizer.transform(expanded_test_X))

    sbert_rf_seed_pred = run_pred(
        sbert_rf_classifier, sbert_featurizer.encode(seed_test_X))
    sbert_rf_expanded_pred = run_pred(
        sbert_rf_classifier, sbert_featurizer.encode(expanded_test_X))

    sbert_svm_seed_pred = run_pred(
        sbert_svm_classifier, sbert_featurizer.encode(seed_test_X))
    sbert_svm_expanded_pred = run_pred(
        sbert_svm_classifier, sbert_featurizer.encode(expanded_test_X))

    print('================tfidf + RF  Model======================')
    print('baseline: ')
    print(classification_report(seed_test_Y, tfidf_rf_seed_pred))
    print(' ')
    print(f'tricking model with {run_name}: ')
    print(classification_report(expanded_test_Y, tfidf_rf_expanded_pred))
    print(' ')

    print('================tfidf + SVM  Model======================')
    print('baseline: ')
    print(classification_report(seed_test_Y, tfidf_svm_seed_pred))
    print(' ')
    print(f'tricking model with {run_name}: ')
    print(classification_report(expanded_test_Y, tfidf_svm_expanded_pred))
    print(' ')

    print('================sbert + RF  Model======================')
    print('baseline: ')
    print(classification_report(seed_test_Y, sbert_rf_seed_pred))
    print(' ')
    print(f'tricking model with {run_name}: ')
    print(classification_report(expanded_test_Y, sbert_rf_expanded_pred))
    print(' ')

    print('================sbert + SVM  Model======================')
    print('baseline: ')
    print(classification_report(seed_test_Y, sbert_svm_seed_pred))
    print(' ')
    print(f'tricking model with {run_name}: ')
    print(classification_report(expanded_test_Y, sbert_svm_expanded_pred))
    print(' ')

    print('================Result Summary=======================')
    tfidf_rf_result = check_basic_assertion(
        seed_test_Y, tfidf_rf_seed_pred, expanded_test_Y, tfidf_rf_expanded_pred, delta=DELTA_DIFF)
    tfidf_svm_result = check_basic_assertion(
        seed_test_Y, tfidf_svm_seed_pred, expanded_test_Y, tfidf_svm_expanded_pred, delta=DELTA_DIFF)
    sbert_rf_result = check_basic_assertion(
        seed_test_Y, sbert_rf_seed_pred, expanded_test_Y, sbert_rf_expanded_pred, delta=DELTA_DIFF)
    sbert_svm_result = check_basic_assertion(
        seed_test_Y, sbert_svm_seed_pred, expanded_test_Y, sbert_svm_expanded_pred, delta=DELTA_DIFF)
    if (tfidf_rf_result and tfidf_svm_result and sbert_rf_result and sbert_svm_result):
        print(f'{run_name} Passing all tests.')
    else:
        print(f'{run_name} Failling basic model tests.')
        if (not tfidf_rf_result):
            print('TFIDF RF Failed.')
        if (not tfidf_svm_result):
            print('TFIDF SVM Failed.')
        if (not sbert_rf_result):
            print('SBERT RF Failed.')
        if (not sbert_svm_result):
            print('SBERT SVM Failed.')
        assert False
