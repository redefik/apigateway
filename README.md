# apigateway
Api Gateway per il front-end Android dell'applicazione a microservizi di supporto alla didattica universitaria.

Si occupa di inoltrare le richieste ai microservizi componenti dell'applicazione, autenticandole e autorizzandole.
Quando l'assolvimento di un compito richiede la partecipazione di due o più microservizi decompone la richiesta del client in più sottorichieste effettuandole se possibile in parallelo.
Il client comunica con l'Api Gateway tramite l'interfaccia REST, i cui endpoint sono documentati in [api](api/)

## Linguaggio
Go
