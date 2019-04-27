from entities.delivery_boy import DeliveryBoy


class DeliveryBoyManager:
    def __init__(self, env, count):
        # self.registry = simpy.Store(env, count)
        self.env = env
        self.dbGenerator = self.delivery_boy_generator(0)

    def getdeliveryboy(self):
        return self.dbGenerator.next()

    def delivery_boy_generator(self, i):
        while True:
            i += 1
            yield DeliveryBoy(self.env, i)
