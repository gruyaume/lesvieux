import { useState } from "react";
import { BlogPost } from "../../../types";
import { Panel, EmptyState, MainTable, ContextualMenu, ConfirmationModal } from "@canonical/react-components";
import { RequiredBlogPostParams, deleteBlogPost } from "../../../queries"
import { UseMutationResult, useMutation, useQueryClient } from "react-query"
import { useCookies } from "react-cookie"
import { formatDate } from "../../../utils";


function BlogEmptyState({ }: {}) {
    return (
        <EmptyState
            image={""}
            title="No Blog Posts Available yet."
        >
            <p>
                There are no blog posts in LesVieux.
            </p>
        </EmptyState>
    );
}

export type ConfirmationModalData = {
    onMouseDownFunc: () => void
    warningText: string
} | null

type TableProps = {
    blogPosts: BlogPost[];
};

export function BlogPostsTable({ blogPosts: rows }: TableProps) {
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    const [confirmationModalData, setConfirmationModalData] = useState<ConfirmationModalData | null>(null);

    const mutationFunc = (mutation: UseMutationResult<any, unknown, RequiredBlogPostParams, unknown>, params: RequiredBlogPostParams) => {
        mutation.mutate(params)
    }
    const queryClient = useQueryClient()
    const deleteMutation = useMutation(deleteBlogPost, {
        onSuccess: () => {
            queryClient.invalidateQueries('blogposts')
        }
    })

    const handleDelete = (id: number) => {
        setConfirmationModalData({
            onMouseDownFunc: () => mutationFunc(deleteMutation, { id: id.toString(), authToken: cookies.user_token }),
            warningText: "Deleting a blog post means it will be completely removed forever. This action cannot be undone.",
        });
    };


    const blogPostsRows = rows.map((blogPost) => {
        const { id, title, author, status, created_at } = blogPost;
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
                    emptyStateMsg={<BlogEmptyState />}
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
                    rows={blogPostsRows}
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
