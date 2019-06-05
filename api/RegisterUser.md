**Register User**
----
  Registers an user with the information provided in the JSON body of the request.

* **URL**

  /users/

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   None
   

* **Data Params**

  `{username: "admin", password:"admin_pass", name:"John", surname:"Doe", fiscalCode:"teacherFiscalCode"}`

* **Success Response:**

  * **Code:** 201 CREATED <br />
    **Content:** `{ user : {
                        name: "John",
                        surname: "Doe",
                        type: "teacher",
                        username:"admin",
                        password:"hashed_password"}
                  }`
 
* **Error Response:**

  * **Code:** 409 CONFLICT <br />
    **Content:** `{ error : "Conflict - Username already exists}`

  OR

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{ error : "Bad request" }`
    
  OR

  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "Unauthorized registration" }`
    
  OR

  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal Server Error" }`