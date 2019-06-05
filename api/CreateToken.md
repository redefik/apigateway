**Create token**
----
    Creates a JWT access token for the given user. The token will be used 
    for authenticating next requests.
* **URL**

  /token/

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   None
   

* **Data Params**

   `{username: "admin", password: "admin_pass"}`
   
* **Success Response:**

  * **Code:** 201 CREATED <br />
    **Set-Cookie**: `token=<JWTtoken>`
 
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "Authentication failed - Wrong username or password"}`

  OR

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{ error : "Bad request" }`
    
  OR

  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Api Gateway - Internal Server Error" }`