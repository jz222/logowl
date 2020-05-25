package templates

const ResetPassword = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />

    <title>Log Owl | Sign Up</title>

    <style type="text/css">
        html {
            font-family: Helvetica, sans-serif;
        }

        body {
            margin: 0;
            padding: 0;
        }

        a {
            margin: 20px 0;
            padding: 15px 25px;
            box-sizing: border-box;
            border-radius: 7px;
            background: #0368ff;
            text-decoration: none;
            font-size: 14px;
            color: white !important;
		}
		
		a:focus,
		a:active {
			color: white;
		}

        img {
            height: 70px;
        }

        h1 {
            margin-top: 20px;
            font-size: 24px;
            color: #333;
        }

        hr {
            height: 1px;
            border: none;
            background: #e5ecf9;
        }

        p {
            line-height: 1.5;
            color: #5B6F8C;
        }

        .wrapper {
            width: 100%;
            padding: 25px;
            border-spacing: 0;
        }

        .box {
            width: 100%;
            max-width: 700px;
            padding: 0;
            border: 1px solid #e5ecf9;
            border-spacing: 0;
            border-radius: 7px;
            background-color: white;
        }

        .logo {
            height: 70px;
            padding: 50px 0 25px 0;
            text-align: center;
        }

        .button {
            margin: 40px 0;
        }

        .greeting {
            margin-bottom: 40px;
        }

        .footer {
            margin: 20px 0;
            text-align: center;
            font-size: 10px;
            color: #0368ff;
        }
    </style>
</head>
<body style="background-color: #f4f8fe;">
<table class="wrapper">
    <tr>
        <td>
            <div class="logo">
                <img src="" alt="logo" />
            </div>

            <table class="box" align="center">
                <tr>
                    <td style="padding: 25px;">

                        <h1>Sign Up</h1>

                        <p>Hi {{.FirstName}},</p>

						<p>
							you requested a new password for your Log Owl account. Please click the link below
							to set a new password.
						</p>

						<p>
							Keep in mind that the link is only valid for the next <b>two hours</b>.
						</p>

                        <p class="button">
                            <a href="{{.URL}}">Reset Password</a>
						</p>

                        <p class="greeting">
                            Best Regards,<br />
                            <b>Log Owl</b>
                        </p>

                        <hr />

                        <div class="footer">Log Owl</div>
                    </td>
                </tr>
            </table>
        </td>
    </tr>
</table>
</body>
</html>
`

const ResetPasswordRaw = `
Hi {{.FirstName}},

you requested a new password for your Log Owl account. Please click the link below to set a new password.

Keep in mind that the link is only valid for the next two hours.

{{.URL}}

Best Regards,
Log Owl
`
