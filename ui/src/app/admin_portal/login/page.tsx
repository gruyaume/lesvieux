"use client"

import { getStatus, adminLogin } from "../../queries"
import { useMutation, useQuery } from "react-query"
import { useState, ChangeEvent, useEffect } from "react"
import { useCookies } from "react-cookie"
import { useRouter } from "next/navigation"
import { useAuth } from "../auth/authContext"
import { statusResponseResult } from "../../types"
import Logo from "../../components/logo"
import { Navigation, Notification, Input, PasswordToggle, Button, Form } from "@canonical/react-components";

export default function LoginPage() {
    const router = useRouter()
    const auth = useAuth()
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    const statusQuery = useQuery<statusResponseResult, Error>({
        queryFn: () => getStatus()
    })
    useEffect(() => {
        if (auth.user) {
            router.push("/admin_portal/users");
        }
    }, [auth.user, router]);

    if (!auth.firstUserCreated && (statusQuery.data && !statusQuery.data.initialized)) {
        router.push("/admin_portal/initialize")
    }
    const mutation = useMutation(adminLogin, {
        onSuccess: (response) => {
            const token = response.token;
            if (token) {
                setErrorText("")
                setCookie('user_token', token, {
                    sameSite: true,
                    secure: true,
                    path: "/admin_portal",
                    expires: new Date(new Date().getTime() + 60 * 60 * 1000),
                })
                router.push('/admin_portal/users')
            } else {
                setErrorText("Failed to retrieve token.")
            }
        },
        onError: (e: Error) => {
            setErrorText(e.message)
        }
    })

    const [email, setEmail] = useState<string>("")
    const [password, setPassword] = useState<string>("")

    const [errorText, setErrorText] = useState<string>("")
    const handleEmailChange = (event: ChangeEvent<HTMLInputElement>) => { setEmail(event.target.value) }
    const handlePasswordChange = (event: ChangeEvent<HTMLInputElement>) => { setPassword(event.target.value) }
    return (
        <>
            <Navigation
                items={[]}
                logo={
                    <Logo >
                    </Logo>
                }
            />
            <div style={{
                display: "flex",
                alignContent: "center",
                justifyContent: "center",
                flexWrap: "wrap",
                height: "93.5vh",
            }}>
                <div className="p-panel" style={{
                    width: "35rem",
                    minWidth: "min-content",
                    minHeight: "min-content",
                }}>
                    <div className="p-panel__content">
                        <div className="u-fixed-width">
                            <Form>
                                <fieldset>
                                    <h2 className="p-panel__title">Login</h2>
                                    <Input
                                        id="InputEmail"
                                        label="Email"
                                        type="text"
                                        onChange={handleEmailChange}
                                    />
                                    <PasswordToggle
                                        id="InputPassword"
                                        label="Password"
                                        onChange={handlePasswordChange}
                                    />
                                    {errorText &&
                                        <Notification
                                            severity="negative"
                                            title="Error"
                                        >
                                            {errorText.split("error: ")}
                                        </Notification>
                                    }
                                    <Button
                                        appearance="positive"
                                        disabled={password.length == 0 || email.length == 0}
                                        onClick={
                                            (event) => {
                                                event.preventDefault();
                                                mutation.mutate({ email: email, password: password })
                                            }
                                        }
                                    >
                                        Log In
                                    </Button>
                                </fieldset>
                            </Form>
                        </div>
                    </div>
                </div >
            </div >
        </>
    )
}