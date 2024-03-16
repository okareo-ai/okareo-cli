from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.ensemble import RandomForestClassifier
from sklearn.svm import SVC
import requests
import json
from sentence_transformers import SentenceTransformer
from sklearn.metrics import precision_recall_fscore_support as score


def read_in_data_points(file_path, max_line=None):
    data = []
    count = 0
    with open(file_path) as file:
        lines = file.readlines()
        for line in lines:
            count += 1
            parts = line.split("\t")
            assert len(parts) == 2
            data.append((parts[1], parts[0]))
            if max_line and count == max_line:
                break
    return data


def run_pred(classifier, features_test):
    predictions = classifier.predict(features_test)
    return list(predictions.flatten())


def get_tfidf_featurizer(train_X):
    vectorizer = TfidfVectorizer()
    vectorizer.fit_transform(train_X)
    return vectorizer


def get_sbert_featurizer():
    model = SentenceTransformer(
        'sentence-transformers/paraphrase-MiniLM-L6-v2')
    return model


def get_svm_classifier(train_X, train_Y):
    svm = SVC(kernel='rbf')  # with rbf kernel
    svm.fit(train_X, train_Y)
    return svm


def get_rf_classifier(train_X, train_Y):
    rf = RandomForestClassifier(n_estimators=50)
    rf.fit(train_X, train_Y)
    return rf


def get_download_scenario_url(scenario_id):
    return f"https://api.okareo.com/v0/scenario_sets_download/{scenario_id}"


def download_scenario_data_points(scenario_id, api_key):
    url = get_download_scenario_url(scenario_id)
    headers = {
        'Content-Type': 'application/json',
        'accept': 'application/x-ndjson',
        'api-key': api_key
    }
    return requests.get(url, headers=headers)


def parse_sceenario_json_data(json_response):
    json_response_data = json_response.split('\n')
    expanded_data_X, expanded_data_Y = [], []
    for line in json_response_data:
        json_data = json.loads(line)
        data_point, label = json_data['input'], json_data['result']
        expanded_data_X.append(data_point)
        expanded_data_Y.append(label)
    assert (len(expanded_data_X) == len(expanded_data_Y))
    return expanded_data_X, expanded_data_Y


def check_basic_assertion(baseline_test_Y, pred_Y, attack_test_Y, attack_pred_Y, delta):
    precision_s1, recall_s1, fscore_s1, _support_s1 = score(
        baseline_test_Y, pred_Y, average='weighted')
    precision_s2, recall_s2, fscore_s2, _support_s2 = score(
        attack_test_Y, attack_pred_Y, average='weighted')
    return (precision_s1 - delta > precision_s2) and (recall_s1 - delta > recall_s2) and (fscore_s1 - delta > fscore_s2)
