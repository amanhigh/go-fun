from locust import HttpUser, task, between, tag
import random
import json
import logging
import string

API_VERSION = "v1"

# https://docs.locust.io/en/stable/writing-a-locustfile.html
class FunAppUser(HttpUser):
    wait_time = between(1, 5)
    student_ids = []

    # Task with Weight 1 which is Least Frequent, Higher the number, more frequent the task
    @task(1)
    @tag('telemetry')
    def get_metrics(self):
        self.client.get("/metrics")

    @task(3)
    @tag('write')
    def create_student(self):
        payload = self._generate_student_payload()
        with self.client.post(f"/{API_VERSION}/student", json=payload, catch_response=True) as response:
            if self._handle_response(response, "create") == 201:
                student_id = response.json().get('id')
                if student_id:
                    self.student_ids.append(student_id)
                else:
                    logging.error(f"Student created but no ID found in response. Response: {response.text}")

    @task(5)
    @tag('read')
    def get_student(self):
        if self.student_ids:
            student_id = random.choice(self.student_ids)
            self.client.get(f"/{API_VERSION}/student/{student_id}")
        else:
            # If no students have been created yet, create one
            self.create_student()
            
    @task(4)
    @tag('read')
    def list_students(self):
        params = {
            "offset": random.randint(1, 50),
            "limit": random.randint(2, 5)
        }
        self.client.get(f"/{API_VERSION}/student", params=params)
    
    @task(3)
    @tag('search')
    def get_student_audit(self):
        if self.student_ids:
            student_id = random.choice(self.student_ids)
            with self.client.get(f"/{API_VERSION}/student/{student_id}/audit", catch_response=True) as response:
                if response.status_code != 200:
                    error_msg = f"Failed to get audit for student {student_id}. Status code: {response.status_code}, Response: {response.text}"
                    logging.error(error_msg)
                    response.failure(error_msg)
        else:
            # If no students have been created yet, create one
            self.create_student()

    @task(2)
    @tag('write')
    def update_student(self):
        if self.student_ids:
            student_id = random.choice(self.student_ids)
            payload = self._generate_student_payload()
            payload["name"] = f"Updated {payload['name']}"
            with self.client.put(f"/{API_VERSION}/student/{student_id}", json=payload, catch_response=True) as response:
                self._handle_response(response, "update", student_id)
        else:
            self.create_student()
    
    @task(1)
    @tag('write')
    def delete_student(self):
        if self.student_ids:
            student_id = random.choice(self.student_ids)
            with self.client.delete(f"/{API_VERSION}/student/{student_id}", catch_response=True) as response:
                if response.status_code == 204:
                    self.student_ids.remove(student_id)
                else:
                    error_msg = f"Failed to delete student {student_id}. Status code: {response.status_code}, Response: {response.text}"
                    logging.error(error_msg)
                    response.failure(error_msg)
        else:
            # If no students have been created yet, create one
            self.create_student()

    @task(3)
    @tag('search')
    def list_students_with_sorting(self):
        sort_by = random.choice(["name", "gender", "age"])
        order = random.choice(["asc", "desc"])
        params = {
            "offset": random.randint(1, 50),
            "limit": random.randint(2, 5),
            "sort_by": sort_by,
            "order": order
        }
        with self.client.get(f"/{API_VERSION}/student", params=params, catch_response=True) as response:
            if response.status_code != 200:
                error_msg = f"Failed to list students with sorting. Status code: {response.status_code}, Response: {response.text}"
                logging.error(error_msg)
                response.failure(error_msg)

    @task(3)
    @tag('search')
    def list_students_with_filtering(self):
        params = self._get_filter_params()
        with self.client.get(f"/{API_VERSION}/student", params=params, catch_response=True) as response:
            self._handle_response(response, "list with filtering")

### Helpers
    def _handle_response(self, response, action, student_id=None):
        if response.status_code not in [200, 201, 204]:
            error_msg = f"Failed to {action} student{' ' + str(student_id) if student_id else ''}. Status code: {response.status_code}, Response: {response.text}"
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
    def _generate_student_payload():
        return {
            "name": FunAppUser._generate_username(),
            "age": FunAppUser._generate_age(),
            "gender": FunAppUser._generate_gender()
        }