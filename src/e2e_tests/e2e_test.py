from confest import db_all_tables, api_url
import requests
import classes

def test_auth_and_create_entry(db_all_tables, api_url):
    input_data_user = {
        "about": "string",
        "email": "string",
        "name": "string",
        "password": "string"
    }

    response_auth = requests.post(f"{api_url}/signup", 
                             json=input_data_user)
    user = db_all_tables.query(classes.User).first()

    assert response_auth.status_code == 201
    assert user.email == input_data_user["email"]

    input_data_entry = {
        "description": "work",
        "time_end": "2018-09-24T13:42:31Z",
        "time_start": "2018-09-23T12:42:31Z"
    }
    headers = {"Cookie": response_auth.headers["Set-Cookie"]}
    response_entry = requests.post(f"{api_url}/entry/create", 
                             json=input_data_entry, headers=headers)
    entry = db_all_tables.query(classes.Entry).first()

    assert response_entry.status_code == 201
    # assert response_entry.json["duration"] == "1h"
    assert response_entry.json()["body"]["duration"] == "25h0m0s"
    assert entry.description == input_data_entry["description"]