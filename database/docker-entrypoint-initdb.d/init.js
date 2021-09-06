db.createUser({
    user: 'mongo_user',
    pwd: 'mongo_pass',
    roles: [
        {
            role: 'readWrite',
            db: 'sw-api'
        }
    ]
})