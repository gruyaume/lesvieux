import { Dispatch, SetStateAction, useState, createContext } from "react"
import { useAuth } from "../../../auth/authContext"
import { useMutation, useQueryClient } from "react-query"
import { changeEmployerAccountPassword, createEmployerUser } from "../../../../queries"
import { Modal, Button, Input, PasswordToggle, Form, Notification } from "@canonical/react-components";
import { useFormik } from "formik";
import * as Yup from "yup";

export type ChangeEmployerPasswordModalData = {
    EmployerId: string
    Accountid: string
    email: string
} | null

export type CreateEmployerAccountModalData = {
    EmployerId: string
} | null

interface ChangePasswordModalProps {
    modalData: ChangeEmployerPasswordModalData
    setModalData: Dispatch<SetStateAction<ChangeEmployerPasswordModalData>>
}

type CreateUserModalProps = {
    modalData: CreateEmployerAccountModalData
    setModalData: Dispatch<SetStateAction<boolean>>;
};


export const ChangePasswordModalContext = createContext<ChangePasswordModalProps>({
    modalData: null,
    setModalData: () => { }
})


export function CreateEmployerUserModal({ modalData, setModalData }: CreateUserModalProps) {
    const auth = useAuth();
    const queryClient = useQueryClient();
    const [errorText, setErrorText] = useState<string>("");

    const mutation = useMutation(createEmployerUser, {
        onSuccess: () => {
            queryClient.invalidateQueries("employer_users");
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
        },
        validationSchema: Yup.object().shape({
            email: Yup.string()
                .email("Invalid email")
                .required("Email is required"),
            password1: Yup.string()
                .matches(
                    /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9@#$%^&+=])(?=.{8,})/,
                    "Password must have 8 or more characters, include at least one capital letter, one lowercase letter, and either a number or a symbol"
                )
                .required("Password is required"),
            password2: Yup.string()
                .oneOf([Yup.ref("password1")], "Passwords must match")
                .required("Please confirm your password"),
        }),
        onSubmit: (values) => {
            mutation.mutate({
                employerId: modalData?.EmployerId ?? "",
                authToken: auth.user ? auth.user.authToken : "",
                email: values.email,
                password: values.password1,
            });
        },
    });

    return (
        <Modal
            title={"Create Employer User"}
            buttonRow={
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
                    {errorText && (
                        <Notification
                            inline
                            borderless
                            severity="negative"
                            title="Error:"
                            style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', textAlign: 'left' }}
                        >
                            {errorText}
                        </Notification>
                    )}
                    <div style={{ display: 'flex', gap: '0.5rem', marginLeft: 'auto' }}>
                        <Button onClick={() => setModalData(false)}>Cancel</Button>
                        <Button
                            appearance="positive"
                            disabled={!formik.dirty || !formik.isValid}
                            onClick={formik.submitForm}
                        >
                            Submit
                        </Button>
                    </div>
                </div>
            }
        >
            <Form>
                <Input
                    type="email"
                    id="email"
                    label="Email"
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

export function ChangeEmployerUserPasswordModal({ modalData, setModalData }: ChangePasswordModalProps) {
    const auth = useAuth();
    const queryClient = useQueryClient();
    const [errorText, setErrorText] = useState<string>("");

    const mutation = useMutation(changeEmployerAccountPassword, {
        onSuccess: () => {
            queryClient.invalidateQueries("employer_users");
            setErrorText("");
            setModalData(null);
        },
        onError: (e: Error) => {
            setErrorText(e.message);
        },
    });

    const formik = useFormik({
        initialValues: {
            password1: "",
            password2: "",
        },
        validationSchema: Yup.object().shape({
            password1: Yup.string()
                .matches(
                    /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9@#$%^&+=])(?=.{8,})/,
                    "Password must have 8 or more characters, include at least one capital letter, one lowercase letter, and either a number or a symbol"
                )
                .required("Password is required"),
            password2: Yup.string()
                .oneOf([Yup.ref("password1")], "Passwords must match")
                .required("Please confirm your password"),
        }),
        onSubmit: (values) => {
            mutation.mutate({
                authToken: auth.user ? auth.user.authToken : "",
                employerId: modalData?.EmployerId ?? "",
                accountId: modalData?.Accountid ?? "",
                password: values.password1,
            });
        },
    });

    return (
        <Modal
            title="Change User Password"
            buttonRow={
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
                    {errorText && (
                        <Notification
                            inline
                            borderless
                            severity="negative"
                            title="Error:"
                            style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', textAlign: 'left' }}
                        >
                            {errorText}
                        </Notification>
                    )}
                    <div style={{ display: 'flex', gap: '0.5rem', marginLeft: 'auto' }}>
                        <Button onClick={() => setModalData(null)}>Cancel</Button>
                        <Button
                            appearance="positive"
                            disabled={!formik.dirty || !formik.isValid}
                            onClick={formik.submitForm}
                        >
                            Submit
                        </Button>
                    </div>
                </div>
            }
        >
            <Form>
                <PasswordToggle
                    id="password1"
                    label="New Password"
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
