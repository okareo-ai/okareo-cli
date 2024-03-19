from sklearn.metrics import classification_report
from okareo import Okareo
from okareo_api_client.models import ScenarioSetCreate, SeedData
from utils import read_in_data_points, run_pred, download_scenario_data_points, parse_sceenario_json_data, check_basic_assertion
import random
import string
import os
from model_utils import TfidfSvmClassificationModel, TfidfRFClassificationModel, SbertSvmClassificationModel, SbertRFClassificationModel

OKAREO_API_KEY = os.environ["OKAREO_API_KEY"]
okareo = Okareo(OKAREO_API_KEY)
DELTA_DIFF = 0.02  # Large delta is more strict and better

file_path = os.path.join(
    os.path.dirname(__file__), "datasets/sentiment/binary_small.txt"
)


def get_scenario_set(
    name,
    number_examples,
    generation_type,
    seed_data
):
    scenario_set_create = ScenarioSetCreate(
        name=name,
        number_examples=number_examples,
        generation_type=generation_type,
        seed_data=seed_data
    )

    return okareo.create_scenario_set(scenario_set_create)


def run_basic_generator(run_name, scenario_type):
    ########################## Read In Data Points ###########################
    data = read_in_data_points(file_path)
    seed_data = [SeedData(input_=x, result=y) for x, y in data]
    print(f'sample size: {len(seed_data)}')

    ########################## Create/Download Scenario Set ##################
    random_string = ''.join(random.choices(string.ascii_letters, k=10))
    max_num_samples = 5

    scenario = get_scenario_set(
        name=f"GitHub {run_name} Scenario Categories Baseline " +
        random_string,
        number_examples=max_num_samples,
        generation_type=scenario_type,
        seed_data=seed_data
    )
    scenario_id = scenario.scenario_id

    train_X, train_Y = zip(*data)

    ### Creating Models ###
    tfidf_svm_classifier = TfidfSvmClassificationModel(
        name="tfidf_svm classification model - " + random_string
    )
    tfidf_svm_classifier.fit(train_X, train_Y)

    tfidf_rf_classifier = TfidfRFClassificationModel(
        name="tfidf_rf classification model - " + random_string
    )
    tfidf_rf_classifier.fit(train_X, train_Y)

    sbert_svm_classifier = SbertSvmClassificationModel(
        name="sbert_svm classification model - " + random_string
    )
    sbert_svm_classifier.fit(train_X, train_Y)

    sbert_rf_classifier = SbertRFClassificationModel(
        name="sbert_rf classification model - " + random_string
    )
    sbert_rf_classifier.fit(train_X, train_Y)

    run_basic_local_test(
        run_name,
        scenario_id,
        data,
        tfidf_rf_classifier,
        tfidf_svm_classifier,
        sbert_rf_classifier,
        sbert_svm_classifier,
    )

    run_basic_okareo_test(
        run_name,
        scenario_id,
        random_string,
        tfidf_rf_classifier,
        tfidf_svm_classifier,
        sbert_rf_classifier,
        sbert_svm_classifier,
    )


def run_basic_okareo_test(
    run_name,
    scenario_id,
    random_string,
    tfidf_rf_classifier,
    tfidf_svm_classifier,
    sbert_rf_classifier,
    sbert_svm_classifier,
):
    print('================tfidf + SVM  Model======================')
    tfidf_svm_classifier_mut = okareo.register_model(
        name=f"GitHub {run_name} : tfidf svm classifier - " + random_string,
        model=tfidf_svm_classifier
    )

    # use the scenario or scenario id to run the test
    tfidf_svm_classifier_eval = tfidf_svm_classifier_mut.run_test(
        scenario=scenario_id,
        name=f"GitHub {run_name} : tfidf svm classifier Eval " + random_string,
        calculate_metrics=True
    )
    print(tfidf_svm_classifier_eval.model_metrics.to_dict())
    print('')

    print('================tfidf + RF  Model======================')
    tfidf_rf_classifier_mut = okareo.register_model(
        name=f"GitHub {run_name} : tfidf rf classifier - " + random_string,
        model=tfidf_rf_classifier
    )

    # use the scenario or scenario id to run the test
    tfidf_rf_classifier_eval = tfidf_rf_classifier_mut.run_test(
        scenario=scenario_id,
        name=f"GitHub {run_name} : tfidf rf classifier Eval " + random_string,
        calculate_metrics=True
    )
    print(tfidf_rf_classifier_eval.model_metrics.to_dict())
    print('')

    print('================sbert + SVM  Model======================')
    sbert_svm_classifier_mut = okareo.register_model(
        name=f"GitHub {run_name} : sbert svm classifier - " + random_string,
        model=sbert_svm_classifier
    )

    # use the scenario or scenario id to run the test
    sbert_svm_classifier_eval = sbert_svm_classifier_mut.run_test(
        scenario=scenario_id,
        name=f"GitHub {run_name} : sbert svm classifier Eval " + random_string,
        calculate_metrics=True
    )
    print(sbert_svm_classifier_eval.model_metrics.to_dict())
    print('')

    print('================sbert + RF  Model======================')
    sbert_rf_classifier_mut = okareo.register_model(
        name=f"GitHub {run_name} : sbert rf classifier - " + random_string,
        model=sbert_rf_classifier
    )

    # use the scenario or scenario id to run the test
    sbert_rf_classifier_eval = sbert_rf_classifier_mut.run_test(
        scenario=scenario_id,
        name=f"GitHub {run_name} : sbert rf classifier Eval " + random_string,
        calculate_metrics=True
    )
    print(sbert_rf_classifier_eval.model_metrics.to_dict())
    print('')


def run_basic_local_test(
    run_name,
    scenario_id,
    data,
    tfidf_rf_classifier,
    tfidf_svm_classifier,
    sbert_rf_classifier,
    sbert_svm_classifier,
):

    response = download_scenario_data_points(scenario_id, OKAREO_API_KEY)
    expanded_data_X, expanded_data_Y = parse_sceenario_json_data(response.text)

    train_X, train_Y = zip(*data)
    seed_test_X, seed_test_Y = list(train_X), list(train_Y)
    expanded_test_X, expanded_test_Y = expanded_data_X, expanded_data_Y

    print(f'generated sample size: {len(expanded_data_X)}')

    ### Running Local Predictions ###
    tfidf_rf_seed_pred = run_pred(tfidf_rf_classifier, seed_test_X)
    tfidf_rf_expanded_pred = run_pred(tfidf_rf_classifier, expanded_test_X)

    tfidf_svm_seed_pred = run_pred(tfidf_svm_classifier, seed_test_X)
    tfidf_svm_expanded_pred = run_pred(tfidf_svm_classifier, expanded_test_X)

    sbert_rf_seed_pred = run_pred(sbert_rf_classifier, seed_test_X)
    sbert_rf_expanded_pred = run_pred(sbert_rf_classifier, expanded_test_X)

    sbert_svm_seed_pred = run_pred(sbert_svm_classifier, seed_test_X)
    sbert_svm_expanded_pred = run_pred(sbert_svm_classifier, expanded_test_X)

    print('=========================================================')
    print('================Running Local Tests======================')
    print('=========================================================')

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
