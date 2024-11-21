import { useState } from "react";
import { JobPost } from "../../types";
import { Panel, EmptyState, MainTable, ContextualMenu, ConfirmationModal } from "@canonical/react-components";
import { RequiredJobPostParams, deleteJobPost } from "../../queries"
import { UseMutationResult, useMutation, useQueryClient } from "react-query"
import { useCookies } from "react-cookie"
import { formatDate } from "../../utils";


function JobEmptyState({ }: {}) {
    return (
        <EmptyState
            image={""}
            title="No Job Posts Available yet."
        >
            <p>
                There are no job posts in LesVieux.
            </p>
        </EmptyState>
    );
}

export type ConfirmationModalData = {
    onMouseDownFunc: () => void
    warningText: string
} | null

type TableProps = {
    jobPosts: JobPost[];
};

export function JobPostsTable({ jobPosts: rows }: TableProps) {
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    const [confirmationModalData, setConfirmationModalData] = useState<ConfirmationModalData | null>(null);

    const mutationFunc = (mutation: UseMutationResult<any, unknown, RequiredJobPostParams, unknown>, params: RequiredJobPostParams) => {
        mutation.mutate(params)
    }
    const queryClient = useQueryClient()
    const deleteMutation = useMutation(deleteJobPost, {
        onSuccess: () => {
            queryClient.invalidateQueries('jobposts')
        }
    })

    const handleDelete = (id: number) => {
        setConfirmationModalData({
            onMouseDownFunc: () => mutationFunc(deleteMutation, { id: id.toString(), authToken: cookies.user_token }),
            warningText: "Deleting a job post means it will be completely removed forever. This action cannot be undone.",
        });
    };


    const jobPostsRows = rows.map((jobPost) => {
        const { id, title, author, status, created_at } = jobPost;
        return {
            sortData: {
                id,
                title,
                author,
                status,
                created_at,
            },
            columns: [
                { content: id.toString() },
                { content: title },
                { content: author },
                { content: status },
                {
                    content: formatDate(created_at),
                },
                {
                    content: (
                        <>
                            <ContextualMenu
                                links={[
                                    {
                                        children: "Delete",
                                        onClick: () => handleDelete(id),
                                    },
                                ]}
                                hasToggleIcon
                                position="right"
                            />
                        </>
                    ),
                    className: "u-align--right has-overflow",
                },
            ],
        };
    });
    return (
        <Panel
            title="All Posts"
            stickyHeader
        >
            <div className="u-fixed-width">
                <MainTable
                    emptyStateMsg={<JobEmptyState />}
                    expanding
                    sortable
                    headers={[
                        {
                            content: "ID",
                            sortKey: "id",
                        },
                        {
                            content: "Title",
                            sortKey: "title",
                        },
                        {
                            content: "Author",
                            sortKey: "author",
                        },
                        {
                            content: "Status",
                            sortKey: "status",
                        },
                        {
                            content: "Created At",
                            sortKey: "created_at",
                        },
                        {
                            content: "Actions",
                            className: "u-align--right has-overflow"
                        }
                    ]}
                    rows={jobPostsRows}
                />
                {confirmationModalData != null && (
                    <ConfirmationModal
                        title="Confirm Action"
                        confirmButtonLabel="Confirm"
                        onConfirm={() => {
                            confirmationModalData?.onMouseDownFunc()
                            setConfirmationModalData(null)
                        }}
                        close={() => setConfirmationModalData(null)}
                    >
                        <p>{confirmationModalData.warningText}</p>
                    </ConfirmationModal>
                )}
            </div>
        </Panel >
    );
}