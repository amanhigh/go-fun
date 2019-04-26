class Dish:
    def __init__(self, id):
        self.id = id

    # Assumption Prep time for dish
    # is same across all restaurants.
    def prep_time(self):
        return 5


class Order:
    def __init__(self, id, restaurant, dish):
        self.id = id
        self.restaurant = restaurant
        self.dish = dish

    def customer_drive_time(self):
        return 2

    def customer_handover_time(self):
        return 2
