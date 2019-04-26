import simpy


class Car:
    def __init__(self, env, name, bcs):
        self.name = name
        self.bcs = bcs
        self.env = env
        self.action = env.process(self.run())

    def run(self):
        while True:
            # Simulate driving to the BCS
            yield env.timeout(3)

            # Arrive at Charging spot and request one
            print('%s arriving at %d' % (self.name, env.now))

            # Once charging station is received
            try:
                with bcs.request() as req:
                    yield req

                    # Change Car
                    print('%s got bcs started charging at %d' % (self.name, self.env.now))
                    charge_duration = 5

                    # Wait for change to finish or handle any interruption
                    yield self.env.process(self.charge(charge_duration))

                    # Charging finished leave now.
                    print('%s Leaving bcs at %d' % (self.name, self.env.now))
            except:
                print('%s Was interrupted. Hope, the battery is full enough ...' % self.name)

            # Enjoy your charged Car
            trip_duration = 20
            yield self.env.timeout(trip_duration)

            # Charge Discharged Again
            print('%s car discharged going for recharge at %d' % (self.name, self.env.now))

    def charge(self, duration):
        # Wait for charge duration
        yield self.env.timeout(duration)


def driver(env, car):
    # Wait max 5 minutes
    yield env.timeout(5)

    # Interrupt Car Process
    car.action.interrupt()


env = simpy.Environment()

# Resource Simulation with two charging Stations
bcs = simpy.Resource(env, capacity=2)
car = None
for i in range(10):
    car = Car(env, "Car %d" % i, bcs)

# Driver Interrupts last Car
env.process(driver(env, car))
env.run(until=50)
