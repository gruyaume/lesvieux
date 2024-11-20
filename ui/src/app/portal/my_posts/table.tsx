import { useState } from "react";
import { JobPost } from "../../types";
import { Panel, Button, EmptyState, MainTable, ContextualMenu, ConfirmationModal, Icon } from "@canonical/react-components";
import { RequiredJobPostParams, deleteMyJobPost, createJobPost } from "../../queries"
import { UseMutationResult, useMutation, useQueryClient } from "react-query"
import { useCookies } from "react-cookie"
import { useRouter } from "next/navigation";
import { formatDate } from "../../utils";


function JobEmptyState({ }: {}) {
    const router = useRouter()
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    const queryClient = useQueryClient()
    const createJobPostMutation = useMutation(createJobPost, {
        onSuccess: (data) => {
            console.log("data", data);
            const newPostId = data.id;
            queryClient.invalidateQueries('jobposts');
            router.push(`/portal/my_posts/draft?id=${newPostId}`);
        }
    });

    const handleCreateClick = () => {
        createJobPostMutation.mutate({
            authToken: cookies.user_token,
        });
    };

    return (
        <EmptyState
            image={""}
            title="No Job Posts Available yet."
        >
            <p>
                There are no job posts in LesVieux. Publish your first post!
            </p>
            <Button
                appearance="positive"
                aria-label="add-job-post"
                onClick={handleCreateClick}>
                Publish a job post.
            </Button>
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
    const router = useRouter();

    const mutationFunc = (mutation: UseMutationResult<any, unknown, RequiredJobPostParams, unknown>, params: RequiredJobPostParams) => {
        mutation.mutate(params)
    }
    const queryClient = useQueryClient()
    const deleteMutation = useMutation(deleteMyJobPost, {
        onSuccess: () => {
            queryClient.invalidateQueries('jobposts')
        }
    })

    const createJobPostMutation = useMutation(createJobPost, {
        onSuccess: (data) => {
            console.log("data", data);
            const newPostId = data.id;
            queryClient.invalidateQueries('jobposts');
            router.push(`/portal/my_posts/draft?id=${newPostId}`);
        }
    });

    const handleCreateClick = () => {
        createJobPostMutation.mutate({
            authToken: cookies.user_token,
        });
    };


    const handleDelete = (id: number) => {
        setConfirmationModalData({
            onMouseDownFunc: () => mutationFunc(deleteMutation, { id: id.toString(), authToken: cookies.user_token }),
            warningText: "Deleting a job post means it will be completely removed forever. This action cannot be undone.",
        });
    };

    const handleEdit = (id: number) => {
        router.push(`/portal/my_posts/draft?id=${id}`);
    };

    const jobPostsRows = rows.map((jobPost) => {
        const { id, title, author, status, created_at } = jobPost;
        return {
            sortData: {
                title,
                status,
                created_at,
            },
            columns: [
                { content: title },
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
                                        children: "Edit",
                                        onClick: () => handleEdit(id),
                                    },
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
            title="My Posts"
            stickyHeader
            controls={
                <>
                    <Button
                        appearance="base"
                        onClick={() => router.push('/')}
                    >
                        <Icon
                            name="external-link" />
                        {' '}Go to job
                    </Button>
                    {rows.length > 0 && (
                        <>
                            <Button
                                appearance="positive"
                                onClick={handleCreateClick}
                            >
                                <Icon
                                    light
                                    name="plus" />
                                {' '}Create
                            </Button>
                        </>
                    )}
                </>
            }
        >
            <div className="u-fixed-width">
                <MainTable
                    emptyStateMsg={<JobEmptyState />}
                    expanding
                    sortable
                    headers={[
                        {
                            content: "Title",
                            sortKey: "title",
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
