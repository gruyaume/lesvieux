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
};

export function UsersTable({ adminUsers }: TableProps) {
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

    const users = adminUsers.map((user) => ({
        ...user,
    }));

    return (
        <Panel
            stickyHeader
            title="Admin Users"
            controls={
                <Button appearance="positive" onClick={() => setIsCreateUserModalOpen(true)}>
                    Create Admin User
                </Button>
            }
        >
            <div className="u-fixed-width">
                <MainTable
                    headers={[
                        { content: "Email" },
                        { content: "Actions", className: "u-align--right has-overflow" },
                    ]}
                    rows={users.map((user) => ({
                        columns: [
                            { content: user.email },
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
