import { useState, useContext } from "react"
import { EmployerEntry } from "../../types"
import { useMutation, useQueryClient } from "react-query"
import { Button, ContextualMenu, MainTable, Panel, ConfirmationModal, Modal, Input, Form } from "@canonical/react-components";
import { useAuth } from "../auth/authContext"
import { deleteEmployer, createEmployer } from "../../queries"


export type ConfirmationModalData = {
    onMouseDownFunc: () => void
    warningText: string
} | null


type TableProps = {
    employers: EmployerEntry[]
}

export function EmployersTable({ employers }: TableProps) {
    const auth = useAuth()
    const [confirmationModalData, setConfirmationModalData] = useState<ConfirmationModalData | null>(null)
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [newEmployerName, setNewEmployerName] = useState("");
    const queryClient = useQueryClient()
    const deleteMutation = useMutation(deleteEmployer, {
        onSuccess: () => queryClient.invalidateQueries('employers')
    })

    const handleDelete = (id: string, email: string) => {
        setConfirmationModalData({
            warningText: `Deleting employer: "${email}". This action cannot be undone.`,
            onMouseDownFunc: () => {
                const authToken = auth.user ? auth.user.authToken : "";
                deleteMutation.mutate({ id: id, authToken });
            }
        });
    };

    const handleCreateEmployer = () => {
        createEmployer({ authToken: auth.user ? auth.user.authToken : "", name: newEmployerName })
            .then(() => {
                queryClient.invalidateQueries("employers");
                setIsCreateModalOpen(false);
                setNewEmployerName("");
            })
            .catch((error) => {
                console.error("Failed to create employer:", error);
            });
    };

    return (
        <Panel
            stickyHeader
            title="Employers"
            controls={
                <Button appearance="positive" onClick={() => setIsCreateModalOpen(true)}>
                    Create Employer
                </Button>
            }
        >
            <div className="u-fixed-width">
                <MainTable
                    headers={[{
                        content: "ID"
                    }, {
                        content: "Name"
                    }, {
                        content: "Actions",
                        className: "u-align--right has-overflow"
                    }]}
                    rows={employers.map(employer => ({
                        columns: [
                            {
                                content: employer.id.toString(),
                            },
                            {
                                content: employer.name,
                            },
                            {
                                content: (
                                    <ContextualMenu
                                        links={[{
                                            children: "Delete Employer",
                                            onClick: () => handleDelete(employer.id.toString(), employer.name)
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
            {isCreateModalOpen && (
                <Modal
                    title="Create Employer"
                    close={() => setIsCreateModalOpen(false)}
                    buttonRow={<>
                        <Button onClick={() => setIsCreateModalOpen(false)}>
                            Cancel
                        </Button>
                        <Button
                            appearance="positive"
                            disabled={newEmployerName.length === 0}
                            onClick={handleCreateEmployer}>
                            Submit
                        </Button>
                    </>}>
                    <Form>
                        <Input
                            label="Employer Name"
                            value={newEmployerName}
                            onChange={(e) => setNewEmployerName(e.target.value)}
                            placeholder="Enter employer name"
                        />
                    </Form>
                </Modal>
            )}

        </Panel>
    )
}