/** @type {import("./$types").Actions} */
export const actions = {
    
    changePassword: async ({request, fetch}) => {

        const data = await request.formData()
        const oldPassword = data.get("old")
        const newPassword = data.get("new")
        const newPasswordConfirm = data.get("newConfirm")

        const req = await fetch(
            "/api/admin/change_password",
            {
                method: "POST",
                body: JSON.stringify({
                    old_password: oldPassword,
                    new_password_0: newPassword,
                    new_password_1: newPasswordConfirm
                })
            }
        )


        if (req.ok) {
            
            return {
                success: true,
                message: "Password changed successfully"
            }

        } else {

            if (req.status == 400) {
                
                return {
                    success: false,
                    message: (await req.text()).split(":")[1]
                }

            } else {
                return {
                    success: false,
                    message: "Server error"
                }
            }

        }

    },
    changeUsername: async ({request}) => {

        

    }

}