import { Dispatch, SetStateAction, useState, ChangeEvent, createContext } from "react"
import { useAuth } from "../../auth/authContext"
import { useMutation, useQueryClient } from "react-query"
import { changePassword, changeMyPassword } from "../../../queries"
import { passwordIsValid } from "../../utils"
import { Modal, Button, Input, PasswordToggle, Form } from "@canonical/react-components";

export type ChangePasswordModalData = {
    id: string
    email: string
} | null

interface ChangePasswordModalProps {
    modalData: ChangePasswordModalData
    setModalData: Dispatch<SetStateAction<ChangePasswordModalData>>
}

export const ChangePasswordModalContext = createContext<ChangePasswordModalProps>({
    modalData: null,
    setModalData: () => { }
})

export function ChangeMyPasswordModal({ modalData, setModalData }: ChangePasswordModalProps) {
    const auth = useAuth()
    const queryClient = useQueryClient()
    const mutation = useMutation(changeMyPassword, {
        onSuccess: () => {
            queryClient.invalidateQueries('users')
            setErrorText("")
            setModalData(null)
        },
        onError: (e: Error) => {
            setErrorText(e.message)
        }
    })
    const [password1, setPassword1] = useState<string>("")
    const [password2, setPassword2] = useState<string>("")
    const passwordsMatch = password1 === password2
    const password1Error = password1 && !passwordIsValid(password1) ? "Password is not valid" : ""
    const password2Error = password2 && !passwordsMatch ? "Passwords do not match" : ""

    const [errorText, setErrorText] = useState<string>("")
    const handlePassword1Change = (event: ChangeEvent<HTMLInputElement>) => { setPassword1(event.target.value) }
    const handlePassword2Change = (event: ChangeEvent<HTMLInputElement>) => { setPassword2(event.target.value) }
    return (
        <Modal
            title="Change My Password"
            buttonRow={<>
                <Button onClick={() => setModalData(null)}>
                    Cancel
                </Button>
                <Button
                    appearance="positive"
                    disabled={!passwordsMatch || !passwordIsValid(password1)}
                    onClick={(event) => { event.preventDefault(); mutation.mutate({ authToken: (auth.user ? auth.user.authToken : ""), password: password1 }) }}>
                    Submit
                </Button>
            </>}>
            <Form>
                <Input
                    id="InputEmail"
                    label="Email"
                    type="text"
                    disabled={true}
                    value={modalData?.email}
                />
                <PasswordToggle
                    help="Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol."
                    id="password1"
                    label="New Password"
                    onChange={handlePassword1Change}
                    error={password1Error}
                />
                <PasswordToggle
                    id="password2"
                    label="Confirm New Password"
                    onChange={handlePassword2Change}
                    error={password2Error}
                />
            </Form>
        </Modal >
    )
}

export function ChangePasswordModal({ modalData, setModalData }: ChangePasswordModalProps) {
    const auth = useAuth()
    const queryClient = useQueryClient()
    const mutation = useMutation(changePassword, {
        onSuccess: () => {
            queryClient.invalidateQueries('users')
            setErrorText("")
            setModalData(null)
        },
        onError: (e: Error) => {
            setErrorText(e.message)
        }
    })
    const [password1, setPassword1] = useState<string>("")
    const [password2, setPassword2] = useState<string>("")
    const passwordsMatch = password1 === password2
    const password1Error = password1 && !passwordIsValid(password1) ? "Password is not valid" : ""
    const password2Error = password2 && !passwordsMatch ? "Passwords do not match" : ""

    const [errorText, setErrorText] = useState<string>("")
    const handlePassword1Change = (event: ChangeEvent<HTMLInputElement>) => { setPassword1(event.target.value) }
    const handlePassword2Change = (event: ChangeEvent<HTMLInputElement>) => { setPassword2(event.target.value) }
    return (
        <Modal
            title="Change User Password"
            buttonRow={<>
                <Button onClick={() => setModalData(null)}>
                    Cancel
                </Button>
                <Button
                    appearance="positive"
                    disabled={!passwordsMatch || !passwordIsValid(password1)}
                    onClick={(event) => { event.preventDefault(); mutation.mutate({ authToken: (auth.user ? auth.user.authToken : ""), id: modalData ? modalData.id : "", password: password1 }) }}>
                    Submit
                </Button>
            </>}>
            <Form>
                <Input
                    id="InputEmail"
                    label="Email"
                    type="text"
                    disabled={true}
                    value={modalData?.email}
                />
                <PasswordToggle
                    help="Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol."
                    id="password1"
                    label="New Password"
                    onChange={handlePassword1Change}
                    error={password1Error}
                />
                <PasswordToggle
                    id="password2"
                    label="Confirm New Password"
                    onChange={handlePassword2Change}
                    error={password2Error}
                />
            </Form>
        </Modal >
    )
}
