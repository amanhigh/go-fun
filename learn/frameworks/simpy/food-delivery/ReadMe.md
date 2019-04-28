# Food Delivery Simulator
## Installation
* Requires Python 2
* pip install simpy

## First Run
> cd food-delivery
> python main.py
 
## Configuring
* All Simulation knobs are in *config.yaml*. 
* Each Parameter is accompanied by documentation on what it does.

Some important config params
* sim.until - Time until which simulation runs to elongate or shorten runs.
* sim.logLevel - CRITICAL shows only summary, INFO Major Steps, DEBUG all steps
* algo.type - Type of algo for computing cost.

## Simpler Run
To See all Steps and shorter run.
Use
restaurant.restaurantCount = 1
restaurant.kitchenCount = 2
delivery.hires=2
sim.unitl=50
sim.loglevel=DEBUG
