import simpy

import entities.delivery_boy as db


class DeliveryBoyManager:
    def __init__(self, env, count):
        self.registry = simpy.Store(env, count)
        self.boy = db.DeliveryBoy(env, 1)
        # for i in range(count):
        #     put = self.registry.put(db.DeliveryBoy(env, i))

    def getdeliveryboy(self):
        return self.boy
