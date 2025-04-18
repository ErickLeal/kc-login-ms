<!DOCTYPE html>
<html>
<head>
    <title>Login - My Company</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f9;
            color: #333;
            text-align: center;
            padding: 50px;
        }
        .login-container {
            max-width: 400px;
            margin: 0 auto;
            padding: 20px;
            background: #fff;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            border: 2px solid #6a0dad; /* Purple border */
        }
        .login-container h1 {
            margin-bottom: 20px;
            color: #6a0dad; /* Purple text */
        }
        .login-container img {
            max-width: 150px;
            margin-bottom: 20px;
        }
        button {
            background-color: #6a0dad; /* Purple button */
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 5px;
            cursor: pointer;
        }
        button:hover {
            background-color: #5a0c9d; /* Darker purple on hover */
        }
    </style>
</head>
<body>
    <div class="login-container">
        <img src="/resources/img/company-logo.png" alt="My Company Logo">
        <h1>Welcome to My Company</h1>
        <form action="${url.loginAction}" method="post">
            <div>
                <label for="username">Username:</label>
                <input type="text" id="username" name="username" required>
            </div>
            <div>
                <label for="password">Password:</label>
                <input type="password" id="password" name="password" required>
            </div>
            <div>
                <button type="submit">Login</button>
            </div>
        </form>
    </div>
</body>
</html>