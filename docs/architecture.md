Description on each of the packages

# cmd
Under the root command, there is only one subcommand named start_server. We use the routers package from it & start the server.

# routes
We have these types of routes: 

`user_auth`: signUp, signIn, signOut

`course`: 
i) get/list -> general authenticated users can do it.
ii) create/update -> admin or moderators can do.
iii) delete -> only admin can do it.

For get,update & delete calls(those work on a specific course uid), we utilize a context middleware for injecting context data.

There is one special route for providing role to a user. It requires the admin access.

# pkg
-`pkg.auth`:

We store these fields in the redis session:
i) authenticated, ii) role, iii) userName, iv) userIP, v) user_agent
There are some getters implemented in the session.go file.



-`pkg.error`

-`pkg.middleware`:

There are 4 types of middlewares in this package.
i) global middleware: For example logger, urlFormatter etc.

ii) common security middleware: To check if a session is valid & authenticated.

iii) access related middleware: To check if the role(in session) is matched.

iv) context middleware: To append additional info to the context.


# handlers
`routes -> handlers -> models`

Routes package will call functions from handlers package. And handlers package will call functions from models package.
We should be strict to this design principle to avoid implementation complexity & circular dependencies.

In general, the handlers are pretty straight-forward. They contain functions like Create, Get, List, Update, Delete etc.

# models
This is the actual working unit. The API & common constants are declared here. Also, it is responsible for all the database calls.


1) `models.db.client`:
A special package only to initialize the db clients thorough init().

2) `models.db`:
DB calls are implemented here. No one calls the DB directly except this package.

3) `models.course`:
Dedicated package for course related methods. Intended to only be called from `handlers/course`.

4) `models.user`:
   Dedicated package for user related methods. Intended to only be called from `handlers/user`.


---
There are some other non-code packages/files worth mentioning.

`docs`: Holds the documentation files.

`hack`: Holds the utility script files.

`.env` file: Export the required ENVs to run this backend server.

`req.http` file: Holds the curl request-response examples for the implemented http routes.
