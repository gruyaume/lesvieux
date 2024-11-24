import { useState } from "react";
import { UserEntry } from "../../../../types";
import { useMutation, useQueryClient } from "react-query";
import {
    Button,
    ContextualMenu,
    MainTable,
    Panel,
    ConfirmationModal,
} from "@canonical/react-components";
import {
    ChangeEmployerPasswordModalData,
    ChangeEmployerUserPasswordModal,
    CreateEmployerAccountModalData,
    CreateEmployerUserModal,
} from "./components";
import { useAuth } from "../../../auth/authContext";
import { deleteEmployerAccount } from "../../../../queries";

export type ConfirmationModalData = {
    onMouseDownFunc: () => void;
    warningText: string;
} | null;

type TableProps = {
    employerUsers: UserEntry[];
    employerId: string;
    employerName: string;
};

export function EmployerUsersTable({ employerUsers, employerId, employerName }: TableProps) {
    const auth = useAuth();
    const queryClient = useQueryClient();

    const [confirmationModalData, setConfirmationModalData] = useState<ConfirmationModalData | null>(null);
    const [changePasswordModalData, setChangePasswordModalData] =
        useState<ChangeEmployerPasswordModalData | null>(null);
    const [isCreateUserModalOpen, setIsCreateUserModalOpen] = useState(false);
    const [createEmployerUserModalData, setCreateEmployerUserModalData] =
        useState<CreateEmployerAccountModalData | null>(null);

    const deleteMutation = useMutation(deleteEmployerAccount, {
        onSuccess: () => queryClient.invalidateQueries(["employer_users"]),
    });

    const handleDelete = (employerId: string, accountId: string, email: string) => {
        setConfirmationModalData({
            warningText: `Deleting user: "${email}". This action cannot be undone.`,
            onMouseDownFunc: () => {
                const authToken = auth.user ? auth.user.authToken : "";
                deleteMutation.mutate({ authToken, employerId, accountId });
            },
        });
    };

    const handleChangePassword = (employerId: string, accountId: string, email: string) => {
        setChangePasswordModalData({
            EmployerId: employerId,
            Accountid: accountId,
            email,
        });
    };

    const handleCreateEmployerUser = () => {
        setCreateEmployerUserModalData({
            EmployerId: employerId,
        });
        setIsCreateUserModalOpen(true);
    };

    const users = employerUsers.map((user) => ({
        ...user,
    }));

    return (
        <Panel
            stickyHeader
            title={`Employers / ${employerName} / Users`}
            controls={
                <Button appearance="positive" onClick={handleCreateEmployerUser}>
                    Create Employer User
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
                                                onClick: () => handleDelete(employerId, user.id.toString(), user.email),
                                            },
                                            {
                                                children: "Change Password",
                                                onClick: () => handleChangePassword(employerId, user.id.toString(), user.email),
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
                <ChangeEmployerUserPasswordModal
                    modalData={changePasswordModalData}
                    setModalData={setChangePasswordModalData}
                />
            )}
            {isCreateUserModalOpen && createEmployerUserModalData && (
                <CreateEmployerUserModal
                    modalData={createEmployerUserModalData}
                    setModalData={setIsCreateUserModalOpen}
                />
            )}
        </Panel>
    );
}
