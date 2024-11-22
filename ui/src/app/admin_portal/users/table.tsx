import { useState } from "react";
import { UserEntry } from "../../types";
import { useMutation, useQueryClient } from "react-query";
import {
    Button,
    ContextualMenu,
    MainTable,
    Panel,
    ConfirmationModal,
} from "@canonical/react-components";
import { ChangeAdminPasswordModalData, ChangeAdminPasswordModal, CreateAccountModalData, CreateUserModal } from "./components";
import { useAuth } from "../auth/authContext";
import { deleteAdminAccount } from "../../queries";

export type ConfirmationModalData = {
    onMouseDownFunc: () => void;
    warningText: string;
} | null;

type TableProps = {
    adminUsers: UserEntry[];
    employerUsers: UserEntry[];
};

export function UsersTable({ adminUsers, employerUsers }: TableProps) {
    const auth = useAuth();
    const queryClient = useQueryClient();

    const [confirmationModalData, setConfirmationModalData] = useState<ConfirmationModalData | null>(null);
    const [changePasswordModalData, setChangePasswordModalData] = useState<ChangeAdminPasswordModalData | null>(null);
    const [isCreateUserModalOpen, setIsCreateUserModalOpen] = useState(false); // Use this state for the modal

    const deleteMutation = useMutation(deleteAdminAccount, {
        onSuccess: () => queryClient.invalidateQueries("admin_users"),
    });

    const handleDelete = (id: string, email: string) => {
        setConfirmationModalData({
            warningText: `Deleting user: "${email}". This action cannot be undone.`,
            onMouseDownFunc: () => {
                const authToken = auth.user ? auth.user.authToken : "";
                deleteMutation.mutate({ id, authToken });
            },
        });
    };

    const handleChangePassword = (id: string, email: string) => {
        setChangePasswordModalData({ id, email });
    };

    const usersWithRoles = [
        ...adminUsers.map((user) => ({ ...user, role: "Admin" })),
        ...employerUsers.map((user) => ({ ...user, role: "Employer" })),
    ];

    return (
        <Panel
            stickyHeader
            title="Users"
            controls={
                <Button appearance="positive" onClick={() => setIsCreateUserModalOpen(true)}>
                    Create User
                </Button>
            }
        >
            <div className="u-fixed-width">
                <MainTable
                    headers={[
                        { content: "Email" },
                        { content: "Role" },
                        { content: "Actions", className: "u-align--right has-overflow" },
                    ]}
                    rows={usersWithRoles.map((user) => ({
                        columns: [
                            { content: user.email },
                            { content: user.role },
                            {
                                content: (
                                    <ContextualMenu
                                        links={[
                                            {
                                                children: "Delete User",
                                                disabled: user.id === 1,
                                                onClick: () => handleDelete(user.id.toString(), user.email),
                                            },
                                            {
                                                children: "Change Password",
                                                onClick: () => handleChangePassword(user.id.toString(), user.email),
                                            },
                                        ]}
                                        hasToggleIcon
                                        position="right"
                                    />
                                ),
                                className: "u-align--right",
                                hasOverflow: true,
                            },
                        ],
                    }))}
                />
            </div>
            {confirmationModalData && (
                <ConfirmationModal
                    title="Confirm Action"
                    confirmButtonLabel="Delete"
                    onConfirm={() => {
                        confirmationModalData?.onMouseDownFunc();
                        setConfirmationModalData(null);
                    }}
                    close={() => setConfirmationModalData(null)}
                >
                    <p>{confirmationModalData?.warningText}</p>
                </ConfirmationModal>
            )}
            {changePasswordModalData && (
                <ChangeAdminPasswordModal
                    modalData={changePasswordModalData}
                    setModalData={setChangePasswordModalData}
                />
            )}
            {isCreateUserModalOpen && (
                <CreateUserModal setModalData={setIsCreateUserModalOpen} />
            )}
        </Panel>
    );
}
