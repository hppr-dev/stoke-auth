# Stoke Client Examples

This directory contains examples clients demonstrating stoke clients.

```mermaid
C4Component
  title Client Example Environment
  Person(endUser, "User", "User")
  
  Boundary(planetExpress, "Planet Express Shipping services") {
    Component(userUI, "User UI", "Javascript", "Interface into the system")

    Boundary(clientsBoundary, "Example Client Services") {
      Component(deliveryRequest, "Delivery Request", "Python (flask)", "Keeps track of what deliveries have been requested")
      Component(companyInventory, "Company Inventory", "Python (django)", "Keeps track of what items the company can deliver")
      Component(shipControl, "Ship Control", "Go (REST)", "Company Inventory")
      Component(cargoHold, "Cargo Hold", "Python (GRPC)", "Hold shipments while in transit")
      Component(engineRoom, "Engine Room", "Go (GRPC)", "Controls engine thrust")
    }
    
    Boundary(extServices, "External Services") {
      System_Ext(stokeServer, "Stoke Server", "Stoke server for requesting/validating user tokens")
      System_Ext(ldap, "LDAP Server", "Holds authentication credentials and external org roles")

      Rel(stokeServer, ldap, "Authenticates against")
      UpdateRelStyle(stokeServer, ldap, $textColor="grey", $offsetX="-40")
    }

    Rel(endUser, userUI, "Loads")
    UpdateRelStyle(endUser, userUI, $textColor="grey", $offsetX="-15", $offsetY="-10")

    Rel(userUI, stokeServer, "Requests Token")
    UpdateRelStyle(userUI, stokeServer, $textColor="grey", $offsetY="35")

    Rel(userUI, shipControl, "Requests commands")
    UpdateRelStyle(userUI, shipControl, $textColor="grey", $offsetX="-85", $offsetY="-20")

    Rel(userUI, deliveryRequest, "Requests delivery")
    UpdateRelStyle(userUI, deliveryRequest, $textColor="grey", $offsetY="-45", $offsetX="-95")

    Rel(userUI, companyInventory, "Inspects/Requests")
    UpdateRelStyle(userUI, companyInventory, $textColor="grey", $offsetX="-40", $offsetY="30")

    Rel(shipControl, engineRoom, "Commands")
    UpdateRelStyle(shipControl, engineRoom, $textColor="grey", $offsetX="-40", $offsetY="10")

    Rel(companyInventory, cargoHold, "Inspects/Places cargo")
    UpdateRelStyle(companyInventory, shipControl, $textColor="grey", $offsetX="-40", $offsetY="5")
    
  }
```

The test environment is a simulation of a space delivery shipment service with the following components:
  * User UI
      * Javascript
      * User Login
      * Request Delivery -- customer role
      * Transfer Cargo -- staff role
      * Inspect ship status -- crew role
      * Inspect Cargo -- crew role
      * Change ship speed -- engineer role
  * Company Inventory
      * Python REST client (django)
      * Keeps track of what items the company can deliver
      * Transfer cargo -- staff role
  * Delivery Request
      * Python REST client (flask)
      * Keeps track of what deliveries have been requested
      * Request delivery -- customer role
  * Ship Control
      * Go REST client
      * Commands or inspects ship
      * Inspect ship status -- crew role
  * Engine
      * Go GRPC client
      * Controls engine thrust
      * Change ship speed -- engineer role
  * Cargo Hold
      * Python GRPC client
      * Hold shipments while in transit
      * Inspect Cargo -- crew role
   
*This is contrived example. The architecture here does not constitute advice*
