import json
from posixpath import split
import numpy as np
import pandas as pd
import string

from datetime import datetime

from sklearn.linear_model import LogisticRegression
from sklearn.feature_extraction.text import CountVectorizer
from nltk.corpus import stopwords
from joblib import dump, load
import math

with open("output.json") as json_file:
    reddit_comments = json.load(json_file)
    
stopwords = set(stopwords.words('english'))
vectorizer = CountVectorizer(token_pattern=r'\b\w+\b')
reddit_data = pd.read_csv('./Reddit_Data.csv')
reddit_data = reddit_data.dropna()
reddit_data.drop_duplicates(subset='clean_comment')


def clean_data(reddit_data):
    data = reddit_data.clean_comment.apply(lambda x: x.split())
    table = str.maketrans('', '', string.punctuation)
    data = data.apply(lambda x: [w.translate(table) for w in x if w.isalpha()])
    data = data.apply(lambda x: [w for w in x if not w in stopwords])
    data = data.apply(lambda x: [w for w in x if len(w) > 1])
    reddit_data["cleaner_column"] = data
    return reddit_data


def split_into_training(reddit_data):

    reddit_data['cleaner_no_list'] = reddit_data.cleaner_column.apply(
        lambda x: ' '.join(x))
    
    reddit_data = reddit_data.reindex(np.random.permutation(reddit_data.index))
    range = math.floor(len(reddit_data)/0.8)
    train = reddit_data[:range]    
    test = reddit_data[range:]
    
    return train, test


def vectorize_training_data(train, test):
    train_matrix = vectorizer.fit_transform(train['cleaner_no_list'])
    test_matrix = vectorizer.transform(test['cleaner_no_list'])

    X_train = train_matrix
    X_test = test_matrix
    y_train = train['category']
    y_test = test['category']

    return X_train, X_test, y_train, y_test


reddit_data = clean_data(reddit_data)
train, test = split_into_training(reddit_data)

X_train, X_test, y_train, y_test = vectorize_training_data(train, test)
lr = LogisticRegression(solver='lbfgs', max_iter=1000)

try:
    lr = load('lr.joblib')
except (OSError, IOError):
    lr.fit(X_train, y_train)
    dump(lr, 'lr.joblib')

predicitions = lr.predict(vectorizer.transform(reddit_comments['Comments']))

negative = [x for x in predicitions if x == 1]
neutral = [x for x in predicitions if x == 0]
positive = [x for x in predicitions if x == -1]

negative_percentage = len(negative)/len(predicitions)
neutral_percentage = len(neutral)/len(predicitions)
positive_percentage = len(positive)/len(predicitions)

percentage_dict = {
    'negative_percent': negative_percentage,
    'neutral_percent': neutral_percentage, 
    'positive_percent': positive_percentage,
    'negative': len(negative),
    'neutral': len(neutral), 
    'positive': len(positive)
}

datetime_string = datetime.now().strftime("%Y-%m-%d")
print(datetime_string)

with open('results.json') as res:
  listObj = json.load(res)
 
listObj[datetime_string] = percentage_dict


with open('results.json', 'w') as json_file:
    json.dump(listObj, json_file, 
                        indent=4,  
                        separators=(',',': '))