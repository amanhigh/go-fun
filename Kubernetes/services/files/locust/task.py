from locust import HttpUser, task, between
import random
import json
import logging
import string

API_VERSION = "v1"

# https://docs.locust.io/en/stable/writing-a-locustfile.html
class FunAppUser(HttpUser):
    # Wait time between tasks in seconds
    wait_time = between(1, 5)
    created_person_ids = []

    # Task with Weight 1 which is Least Frequent, Higher the number, more frequent the task
    @task(1)
    def get_metrics(self):
        self.client.get("/metrics")

    @task(3)
    def create_person(self):
        payload = {
            "name": f"User {random.randint(1, 1000)}",
            "age": random.randint(18, 80),
            "gender": random.choice(["MALE", "FEMALE"])
        }
        with self.client.post(f"/{API_VERSION}/person", json=payload, catch_response=True) as response:
            if response.status_code == 201:
                person_id = response.json().get('id')
                if person_id:
                    self.created_person_ids.append(person_id)
                else:
                    error_msg = f"Person created but no ID found in response. Response: {response.text}"
                    logging.error(error_msg)
                    response.failure(error_msg)
            else:
                error_msg = f"Failed to create person. Status code: {response.status_code}, Response: {response.text}"
                logging.error(error_msg)
                response.failure(error_msg)

    @task(5)
    def get_person(self):
        if self.created_person_ids:
            person_id = random.choice(self.created_person_ids)
            self.client.get(f"/{API_VERSION}/person/{person_id}")
        else:
            # If no persons have been created yet, create one
            self.create_person()
            
    @task(4)
    def list_persons(self):
        limit = random.randint(2, 5)
        params = {
            "offset": 0,
            "limit": limit
        }
        self.client.get(f"/{API_VERSION}/person", params=params)
    
    @task(3)
    def get_person_audit(self):
        if self.created_person_ids:
            person_id = random.choice(self.created_person_ids)
            with self.client.get(f"/{API_VERSION}/person/{person_id}/audit", catch_response=True) as response:
                if response.status_code != 200:
                    error_msg = f"Failed to get audit for person {person_id}. Status code: {response.status_code}, Response: {response.text}"
                    logging.error(error_msg)
                    response.failure(error_msg)
        else:
            # If no persons have been created yet, create one
            self.create_person()

    @task(2)
    def update_person(self):
        if self.created_person_ids:
            person_id = random.choice(self.created_person_ids)
            payload = {
                "name": f"Updated Person {random.randint(1000, 9999)}",
                "age": random.randint(18, 80),
                "gender": random.choice(["MALE", "FEMALE"])
            }
            with self.client.put(f"/{API_VERSION}/person/{person_id}", json=payload, catch_response=True) as response:
                if response.status_code != 200:
                    error_msg = f"Failed to update person {person_id}. Status code: {response.status_code}, Response: {response.text}"
                    logging.error(error_msg)
                    response.failure(error_msg)
        else:
            # If no persons have been created yet, create one
            self.create_person()
    
    @task(1)
    def delete_person(self):
        if self.created_person_ids:
            person_id = random.choice(self.created_person_ids)
            with self.client.delete(f"/{API_VERSION}/person/{person_id}", catch_response=True) as response:
                if response.status_code == 204:
                    self.created_person_ids.remove(person_id)
                else:
                    error_msg = f"Failed to delete person {person_id}. Status code: {response.status_code}, Response: {response.text}"
                    logging.error(error_msg)
                    response.failure(error_msg)
        else:
            # If no persons have been created yet, create one
            self.create_person()

    @task(3)
    def list_persons_with_sorting(self):
        sort_by = random.choice(["name", "gender", "age"])
        order = random.choice(["asc", "desc"])
        limit = random.randint(2, 5)
        params = {
            "offset": 0,
            "limit": limit,
            "sort_by": sort_by,
            "order": order
        }
        with self.client.get(f"/{API_VERSION}/person", params=params, catch_response=True) as response:
            if response.status_code != 200:
                error_msg = f"Failed to list persons with sorting. Status code: {response.status_code}, Response: {response.text}"
                logging.error(error_msg)
                response.failure(error_msg)

    @task(3)
    def list_persons_with_filtering(self):
        # Randomly choose which filter to apply
        filter_type = random.choice(["name", "gender", "age"])
        
        params = {
            "offset": 0,
            "limit": random.randint(2, 5)
        }

        if filter_type == "name":
            # Generate a name similar to how we create users
            name = f"User {random.randint(1, 1000)}"
            params["name"] = name
        elif filter_type == "gender":
            params["gender"] = random.choice(["MALE", "FEMALE"])
        else:  # age
            params["age"] = random.randint(18, 80)

        with self.client.get(f"/{API_VERSION}/person", params=params, catch_response=True) as response:
            if response.status_code != 200:            
                error_msg = f"Failed to list persons with filtering. Status code: {response.status_code}, Response: {response.text}"
                logging.error(error_msg)
                response.failure(error_msg)