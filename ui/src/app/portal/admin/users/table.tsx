import { useState, useContext } from "react"
import { AsideContext } from "../../aside"
import { UserEntry } from "../../../types"
import { useMutation, useQueryClient } from "react-query"
import { Button, ContextualMenu, MainTable, Panel, ConfirmationModal } from "@canonical/react-components";
import { ChangePasswordModalData, ChangePasswordModal } from "./components"
import { useAuth } from "../../auth/authContext"
import { deleteUser } from "../../../queries"


export type ConfirmationModalData = {
    onMouseDownFunc: () => void
    warningText: string
} | null


type TableProps = {
    users: UserEntry[]
}

export function UsersTable({ users }: TableProps) {
    const auth = useAuth()
    const { isOpen: isAsideOpen, setIsOpen: setAsideIsOpen } = useContext(AsideContext)
    const asideContext = useContext(AsideContext)
    const [confirmationModalData, setConfirmationModalData] = useState<ConfirmationModalData | null>(null)
    const [changePasswordModalData, setChangePasswordModalData] = useState<ChangePasswordModalData>(null)
    const queryClient = useQueryClient()
    const deleteMutation = useMutation(deleteUser, {
        onSuccess: () => queryClient.invalidateQueries('users')
    })
    const handleDelete = (id: string, username: string) => {
        setConfirmationModalData({
            warningText: `Deleting user: "${username}". This action cannot be undone.`,
            onMouseDownFunc: () => {
                const authToken = auth.user ? auth.user.authToken : "";
                deleteMutation.mutate({ id: id, authToken });
            }
        });
    };
    const handleChangePassword = (id: string, username: string) => {
        setChangePasswordModalData({ "id": id, "username": username })
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
                        content: "Username"
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
                                content: user.username,
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
                                            onClick: () => handleDelete(user.id.toString(), user.username)
                                        }, {
                                            children: "Change Password",
                                            onClick: () => handleChangePassword(user.id.toString(), user.username)
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
                <ChangePasswordModal
                    modalData={changePasswordModalData}
                    setModalData={setChangePasswordModalData}
                />
            )}
        </Panel>
    )
}
