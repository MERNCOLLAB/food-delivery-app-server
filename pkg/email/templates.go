package email

const ResetPasswordTemplate = `
    <div style="font-family: Arial, sans-serif; line-height:1.6; color: #545454;">
        <h1 style="font-size: 24px; font-weight: bold; color:#d417ff">
            Food Delivery App Password Reset
        </h1>
        <p style="font-size:16px; margin-bottom:10px">
            Your verification code is
            <strong style="font-size: 18px; color:#d417ff;">%s</strong>
        </p>
        <p style="font-size: 14px; margin-bottom: 20px;">
            The code expires in <strong>5 minutes</strong> after this email was sent.
        </p>
        <p style="font-size: 14px;">
            Enter the code in the reset password section of the app to reset your password.
        </p>
        <hr style="border: 0; border-top: 1px solid #ccc; margin: 20px 0;">
        <p style="font-size: 12px; color: #999;">
            If you did not request a password reset, please ignore this email.
        </p>
    </div>`

const SignUpFormTemplate = `
    <div style="font-family: Arial, sans-serif; line-height:1.6; color: #545454;">
        <h1 style="font-size: 24px; font-weight: bold; color:#d417ff">
            Food Delivery App Admin Invitation
        </h1>
        <p style="font-size:16px; margin-bottom:10px">
            You have been invited to sign up as %s.<br>
            Please use the following link to complete your registration:
        </p>
        <p>
            <a href="%s" style="font-size: 18px; color:#d417ff;">%s</a>
        </p>
        <p style="font-size: 14px;">
            This invitation link expires 12 hours after this email is sent.
        </p>
        <p style="font-size: 14px;">
            If you did not expect this invitation, you can ignore this email.
        </p>
    </div>`

const WelcomeWithPasswordTemplate = `
    <div style="font-family: Arial, sans-serif; line-height:1.6; color: #545454;">
        <h1 style="font-size: 24px; font-weight: bold; color:#d417ff">
            Welcome to Food Delivery App!
        </h1>
        <p style="font-size:16px; margin-bottom:10px">
            Your account has been created.<br>
            Here is your temporary password:
            <strong style="font-size: 18px; color:#d417ff;">%s</strong>
        </p>
        <p style="font-size: 14px; margin-bottom: 20px;">
            Please use this password to log in for the first time. You can change your password after logging in.
        </p>
        <hr style="border: 0; border-top: 1px solid #ccc; margin: 20px 0;">
        <p style="font-size: 12px; color: #999;">
            If you did not expect this email, please ignore it.
        </p>
    </div>`
