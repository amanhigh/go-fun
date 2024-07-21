from locust import HttpUser, task, between, tag
import random
import json
import logging
import string

API_VERSION = "v1"

# https://docs.locust.io/en/stable/writing-a-locustfile.html
class FunAppUser(HttpUser):
    wait_time = between(1, 5)
    person_ids = []

    # Task with Weight 1 which is Least Frequent, Higher the number, more frequent the task
    @task(1)
    @tag('telemetry')
    def get_metrics(self):
        self.client.get("/metrics")

    @task(3)
    @tag('write')
    def create_person(self):
        payload = self._generate_person_payload()
        with self.client.post(f"/{API_VERSION}/person", json=payload, catch_response=True) as response:
            if self._handle_response(response, "create") == 201:
                person_id = response.json().get('id')
                if person_id:
                    self.person_ids.append(person_id)
                else:
                    logging.error(f"Person created but no ID found in response. Response: {response.text}")

    @task(5)
    @tag('read')
    def get_person(self):
        if self.person_ids:
            person_id = random.choice(self.person_ids)
            self.client.get(f"/{API_VERSION}/person/{person_id}")
        else:
            # If no persons have been created yet, create one
            self.create_person()
            
    @task(4)
    @tag('read')
    def list_persons(self):
        params = {
            "offset": random.randint(1, 50),
            "limit": random.randint(2, 5)
        }
        self.client.get(f"/{API_VERSION}/person", params=params)
    
    @task(3)
    @tag('search')
    def get_person_audit(self):
        if self.person_ids:
            person_id = random.choice(self.person_ids)
            with self.client.get(f"/{API_VERSION}/person/{person_id}/audit", catch_response=True) as response:
                if response.status_code != 200:
                    error_msg = f"Failed to get audit for person {person_id}. Status code: {response.status_code}, Response: {response.text}"
                    logging.error(error_msg)
                    response.failure(error_msg)
        else:
            # If no persons have been created yet, create one
            self.create_person()

    @task(2)
    @tag('write')
    def update_person(self):
        if self.person_ids:
            person_id = random.choice(self.person_ids)
            payload = self._generate_person_payload()
            payload["name"] = f"Updated {payload['name']}"
            with self.client.put(f"/{API_VERSION}/person/{person_id}", json=payload, catch_response=True) as response:
                self._handle_response(response, "update", person_id)
        else:
            self.create_person()
    
    @task(1)
    @tag('write')
    def delete_person(self):
        if self.person_ids:
            person_id = random.choice(self.person_ids)
            with self.client.delete(f"/{API_VERSION}/person/{person_id}", catch_response=True) as response:
                if response.status_code == 204:
                    self.person_ids.remove(person_id)
                else:
                    error_msg = f"Failed to delete person {person_id}. Status code: {response.status_code}, Response: {response.text}"
                    logging.error(error_msg)
                    response.failure(error_msg)
        else:
            # If no persons have been created yet, create one
            self.create_person()

    @task(3)
    @tag('search')
    def list_persons_with_sorting(self):
        sort_by = random.choice(["name", "gender", "age"])
        order = random.choice(["asc", "desc"])
        params = {
            "offset": random.randint(1, 50),
            "limit": random.randint(2, 5),
            "sort_by": sort_by,
            "order": order
        }
        with self.client.get(f"/{API_VERSION}/person", params=params, catch_response=True) as response:
            if response.status_code != 200:
                error_msg = f"Failed to list persons with sorting. Status code: {response.status_code}, Response: {response.text}"
                logging.error(error_msg)
                response.failure(error_msg)

    @task(3)
    @tag('search')
    def list_persons_with_filtering(self):
        params = self._get_filter_params()
        with self.client.get(f"/{API_VERSION}/person", params=params, catch_response=True) as response:
            self._handle_response(response, "list with filtering")

### Helpers
    def _handle_response(self, response, action, person_id=None):
        if response.status_code not in [200, 201, 204]:
            error_msg = f"Failed to {action} person{' ' + str(person_id) if person_id else ''}. Status code: {response.status_code}, Response: {response.text}"
            logging.error(error_msg)
            response.failure(error_msg)
        return response.status_code

    def _get_filter_params(self):
        filter_type = random.choice(["name", "gender", "age"])
        params = {
            "offset": 0,
            "limit": random.randint(2, 5)
        }
        if filter_type == "name":
            params["name"] = self._generate_username()
        elif filter_type == "gender":
            params["gender"] = self._generate_gender()
        else:  # age
            params["age"] = self._generate_age()
        return params

    @staticmethod
    def _generate_username():
        return f"User {random.randint(1, 1000)}"

    @staticmethod
    def _generate_age():
        return random.randint(18, 80)

    @staticmethod
    def _generate_gender():
        return random.choice(["MALE", "FEMALE"])

    @staticmethod
    def _generate_person_payload():
        return {
            "name": FunAppUser._generate_username(),
            "age": FunAppUser._generate_age(),
            "gender": FunAppUser._generate_gender()
        }