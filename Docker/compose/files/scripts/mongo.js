//Switch to Admin Schema
use admin;

//Create Admin:
 db.createUser(
      {
        user: "admin",
            pwd: "changeme",
            roles: [ { role: "root", db: "admin" } ]
       }
  );

//Show Users
// db.getUsers()

 //Replication Status
 rs.status();

 //Initiate Replicas
 // rs.initiate({
 //     _id: "shard1",
 //     version: 1,
 //     members: [
 //         { _id: 0, host: "mongo1:27017" },
 //         { _id: 1, host: "mongo2:27017" },
 //         { _id: 2, host: "mongo3:27017" }
 //     ]
 // });

 //Member Management
 // rs.add( { host: "mongo1:27017"} ); //Add Member (only on primary as admin)
 // rs.remove("mongo1:27017"); // Remove Member (only on primary as admin)