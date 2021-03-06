import logging
import random

from entities.restaurant import Restaurant
from models.order import Dish


class RestaurantManager:
    def __init__(self, env, config, xy_generator):
        self.xy_generator = xy_generator
        self.restaurantCount = config['restaurantCount']
        self.dishCount = config['dish']['count']

        self.dishMap = self.decide_dishes(config['dish'])
        self.restaurantMap = self.setup_restaurants(env, config)

    def decide_dishes(self, dishConfig):
        dishCount = dishConfig['count']
        minCookTime = dishConfig['cook']['minTime']
        maxCookTime = dishConfig['cook']['maxTime']
        dishMap = {}
        for i in range(dishCount):
            dishMap[i + 1] = Dish(i + 1, random.randint(minCookTime, maxCookTime))
        logging.critical("Deciding %d dishes to be cooked between %d - %d" % (dishCount, minCookTime, maxCookTime))
        return dishMap

    def setup_restaurants(self, env, config):
        count = config['restaurantCount']
        kitchenCount = config['kitchenCount']

        restaurantMap = {}
        logging.critical("Built %d Restaurants with %d Kitchens" % (count, kitchenCount))
        for i in range(count):
            x, y = self.xy_generator.next()
            restaurantMap[i + 1] = Restaurant(env, i + 1, kitchenCount, x, y)
        return restaurantMap

    def get_random_restaurant(self):
        randint = random.randint(1, self.restaurantCount)
        return self.restaurantMap[randint]

    def get_random_dish(self):
        randint = random.randint(1, self.dishCount)
        return self.dishMap[randint]
