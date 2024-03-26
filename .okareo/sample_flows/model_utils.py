from okareo.model_under_test import CustomModel
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.ensemble import RandomForestClassifier
from sklearn.svm import SVC
from sentence_transformers import SentenceTransformer


class TfidfSvmClassificationModel(CustomModel):
    def __init__(self, name):
        self.name = name
        self.featurizer = TfidfVectorizer()
        self.model = SVC(kernel='rbf')  # with rbf kernel

    def fit(self, train_X, train_Y):
        self.featurizer.fit(train_X)
        features_train = self.featurizer.transform(train_X)
        self.model.fit(features_train, train_Y)

    def predict(self, test_X):
        features_test = self.featurizer.transform(test_X)
        return self.model.predict(features_test)

    def invoke(self, input: str):
        pred = self.predict([input])[0]
        return pred, {"label": pred}


class TfidfRFClassificationModel(CustomModel):
    def __init__(self, name):
        self.name = name
        self.featurizer = TfidfVectorizer()
        self.model = RandomForestClassifier(
            n_estimators=50)  # 50 decision trees

    def fit(self, train_X, train_Y):
        self.featurizer.fit(train_X)
        features_train = self.featurizer.transform(train_X)
        self.model.fit(features_train, train_Y)

    def predict(self, test_X):
        features_test = self.featurizer.transform(test_X)
        return self.model.predict(features_test)

    def invoke(self, input: str):
        pred = self.predict([input])[0]
        return pred, {"label": pred}


class SbertSvmClassificationModel(CustomModel):
    def __init__(self, name):
        self.name = name
        self.featurizer = SentenceTransformer(
            'sentence-transformers/paraphrase-MiniLM-L6-v2')
        self.model = SVC(kernel='rbf')  # with rbf kernel

    def fit(self, train_X, train_Y):
        features_train = self.featurizer.encode(train_X)
        self.model.fit(features_train, train_Y)

    def predict(self, test_X):
        features_test = self.featurizer.encode(test_X)
        return self.model.predict(features_test)

    def invoke(self, input: str):
        pred = self.predict([input])[0]
        return pred, {"label": pred}


class SbertRFClassificationModel(CustomModel):
    def __init__(self, name):
        self.name = name
        self.featurizer = SentenceTransformer(
            'sentence-transformers/paraphrase-MiniLM-L6-v2')
        self.model = RandomForestClassifier(
            n_estimators=50)  # 50 decision trees

    def fit(self, train_X, train_Y):
        features_train = self.featurizer.encode(train_X)
        self.model.fit(features_train, train_Y)

    def predict(self, test_X):
        features_test = self.featurizer.encode(test_X)
        return self.model.predict(features_test)

    def invoke(self, input: str):
        pred = self.predict([input])[0]
        return pred, {"label": pred}
