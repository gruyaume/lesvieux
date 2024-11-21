import { useState, useContext } from "react"
import { AsideContext } from "../aside"
import { UserEntry } from "../../types"
import { useMutation, useQueryClient } from "react-query"
import { Button, ContextualMenu, MainTable, Panel, ConfirmationModal } from "@canonical/react-components";
import { ChangeAdminPasswordModalData, ChangeAdminPasswordModal } from "./components"
import { useAuth } from "../auth/authContext"
import { deleteAdminAccount } from "../../queries"


export type ConfirmationModalData = {
    onMouseDownFunc: () => void
    warningText: string
} | null


type TableProps = {
    users: UserEntry[]
}

export function AdminUsersTable({ users }: TableProps) {
    const auth = useAuth()
    const { isOpen: isAsideOpen, setIsOpen: setAsideIsOpen } = useContext(AsideContext)
    const asideContext = useContext(AsideContext)
    const [confirmationModalData, setConfirmationModalData] = useState<ConfirmationModalData | null>(null)
    const [changePasswordModalData, setChangePasswordModalData] = useState<ChangeAdminPasswordModalData>(null)
    const queryClient = useQueryClient()
    const deleteMutation = useMutation(deleteAdminAccount, {
        onSuccess: () => queryClient.invalidateQueries('users')
    })
    const handleDelete = (id: string, email: string) => {
        setConfirmationModalData({
            warningText: `Deleting user: "${email}". This action cannot be undone.`,
            onMouseDownFunc: () => {
                const authToken = auth.user ? auth.user.authToken : "";
                deleteMutation.mutate({ id: id, authToken });
            }
        });
    };
    const handleChangePassword = (id: string, email: string) => {
        setChangePasswordModalData({ "id": id, "email": email })
    }

    return (
        <Panel
            stickyHeader
            title="Users"
            controls={users.length > 0 &&
                <Button appearance="positive" onClick={() => { asideContext.setExtraData(null); setAsideIsOpen(true) }}>Create New User</Button>
            }
        >
            <div className="u-fixed-width">
                <MainTable
                    headers={[{
                        content: "ID"
                    }, {
                        content: "Email"
                    }, {
                        content: "Role"
                    }, {
                        content: "Actions",
                        className: "u-align--right has-overflow"
                    }]}
                    rows={users.map(user => ({
                        columns: [
                            {
                                content: user.id.toString(),
                            },
                            {
                                content: user.email,
                            },
                            {
                                content: user.role === 1 ? "Admin" : "User"
                            },
                            {
                                content: (
                                    <ContextualMenu
                                        links={[{
                                            children: "Delete User",
                                            disabled: user.id === 1,
                                            onClick: () => handleDelete(user.id.toString(), user.email)
                                        }, {
                                            children: "Change Password",
                                            onClick: () => handleChangePassword(user.id.toString(), user.email)
                                        }]}
                                        hasToggleIcon
                                        position="right"
                                    />
                                ),
                                className: "u-align--right",
                                hasOverflow: true
                            }
                        ]
                    }))}
                />
            </div>
            {confirmationModalData && (
                <ConfirmationModal
                    title="Confirm Action"
                    confirmButtonLabel="Delete"
                    onConfirm={() => {
                        confirmationModalData?.onMouseDownFunc()
                        setConfirmationModalData(null)
                    }}
                    close={() => setConfirmationModalData(null)}>
                    <p>{confirmationModalData?.warningText}</p>
                </ConfirmationModal>
            )}
            {changePasswordModalData && (
                <ChangeAdminPasswordModal
                    modalData={changePasswordModalData}
                    setModalData={setChangePasswordModalData}
                />
            )}
        </Panel>
    )
}