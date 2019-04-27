import logging

from entities.restaurant import Restaurant


class RestaurantManager:
    def __init__(self, env, config):
        count = config['restaurantCount']
        kitchenCount = config['kitchenCount']
        self.restaurantMap = {}
        logging.info("Built %d Restaurants with %d Kitchens" % (count, kitchenCount))
        for i in range(count):
            self.restaurantMap[i] = Restaurant(env, i, kitchenCount)
