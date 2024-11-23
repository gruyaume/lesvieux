import { Dispatch, SetStateAction, useState, ChangeEvent, createContext } from "react"
import { useAuth } from "../auth/authContext"
import { useMutation, useQueryClient } from "react-query"
import { changeAdminAccountPassword, changeMyAdminAccountPassword, createAdminUser } from "../../queries"
import { passwordIsValid } from "../../utils"
import { Modal, Button, Input, PasswordToggle, Form, Select } from "@canonical/react-components";
import { useFormik } from "formik";
import * as Yup from "yup";

const validationSchema = Yup.object().shape({
    email: Yup.string().email("Invalid email").required("Email is required"),
    role: Yup.string().required("Role is required"),
    password1: Yup.string()
        .min(8, "Password must be at least 8 characters")
        .required("Password is required"),
    password2: Yup.string()
        .oneOf([Yup.ref("password1")], "Passwords must match")
        .required("Please confirm your password"),
});

export type ChangeAdminPasswordModalData = {
    id: string
    email: string
} | null

export type CreateAccountModalData = {
} | null

interface ChangePasswordModalProps {
    modalData: ChangeAdminPasswordModalData
    setModalData: Dispatch<SetStateAction<ChangeAdminPasswordModalData>>
}

type CreateUserModalProps = {
    setModalData: Dispatch<SetStateAction<boolean>>;
};


export const ChangePasswordModalContext = createContext<ChangePasswordModalProps>({
    modalData: null,
    setModalData: () => { }
})

export function ChangeMyPasswordModal({ modalData, setModalData }: ChangePasswordModalProps) {
    const auth = useAuth()
    const queryClient = useQueryClient()
    const mutation = useMutation(changeMyAdminAccountPassword, {
        onSuccess: () => {
            queryClient.invalidateQueries('admin_users')
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

export function CreateUserModal({ setModalData }: CreateUserModalProps) {
    const auth = useAuth();
    const queryClient = useQueryClient();

    const [errorText, setErrorText] = useState<string>("");

    const mutation = useMutation(createAdminUser, {
        onSuccess: () => {
            queryClient.invalidateQueries("admin_users");
            setErrorText("");
            setModalData(false);
        },
        onError: (e: Error) => {
            setErrorText(e.message);
        },
    });

    const formik = useFormik({
        initialValues: {
            email: "",
            password1: "",
            password2: "",
            role: "0",
        },
        validationSchema,
        onSubmit: (values) => {
            mutation.mutate({
                authToken: auth.user ? auth.user.authToken : "",
                email: values.email,
                password: values.password1,
            });
        },
    });

    return (
        <Modal
            title={"Create Admin User"}
            buttonRow={
                <>
                    <Button onClick={() => setModalData(false)}>Cancel</Button>
                    <Button
                        appearance="positive"
                        onClick={(event) => {
                            event.preventDefault();
                            formik.handleSubmit();
                        }}
                    >
                        Submit
                    </Button>
                </>
            }
        >
            <Form>
                <Input
                    type="text"
                    id="email"
                    label="Email"
                    placeholder="example@lesvieux.ca"
                    required
                    {...formik.getFieldProps("email")}
                    error={formik.touched.email ? formik.errors.email : null}
                />
                <PasswordToggle
                    id="password1"
                    label="Password"
                    required
                    {...formik.getFieldProps("password1")}
                    error={formik.touched.password1 ? formik.errors.password1 : null}
                />
                <PasswordToggle
                    id="password2"
                    label="Confirm Password"
                    required
                    {...formik.getFieldProps("password2")}
                    error={formik.touched.password2 ? formik.errors.password2 : null}
                />
            </Form>
        </Modal>
    );
}

export function ChangeAdminPasswordModal({ modalData, setModalData }: ChangePasswordModalProps) {
    const auth = useAuth()
    const queryClient = useQueryClient()
    const mutation = useMutation(changeAdminAccountPassword, {
        onSuccess: () => {
            queryClient.invalidateQueries('admin_users')
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