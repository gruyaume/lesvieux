import { useState } from "react";
import { BlogPost } from "../../types";
import { Panel, Button, EmptyState, MainTable, ContextualMenu, ConfirmationModal, Icon } from "@canonical/react-components";
import { RequiredBlogPostParams, deleteMyBlogPost, createBlogPost } from "../../queries"
import { UseMutationResult, useMutation, useQueryClient } from "react-query"
import { useCookies } from "react-cookie"
import { useRouter } from "next/navigation";
import { formatDate } from "../../utils";


function BlogEmptyState({ }: {}) {
    const router = useRouter()
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    const queryClient = useQueryClient()
    const createBlogPostMutation = useMutation(createBlogPost, {
        onSuccess: (data) => {
            console.log("data", data);
            const newPostId = data.id;
            queryClient.invalidateQueries('blogposts');
            router.push(`/portal/my_posts/draft?id=${newPostId}`);
        }
    });

    const handleCreateClick = () => {
        createBlogPostMutation.mutate({
            authToken: cookies.user_token,
        });
    };

    return (
        <EmptyState
            image={""}
            title="No Blog Posts Available yet."
        >
            <p>
                There are no blog posts in LesVieux. Publish your first post!
            </p>
            <Button
                appearance="positive"
                aria-label="add-blog-post"
                onClick={handleCreateClick}>
                Publish a blog post.
            </Button>
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
    const router = useRouter();

    const mutationFunc = (mutation: UseMutationResult<any, unknown, RequiredBlogPostParams, unknown>, params: RequiredBlogPostParams) => {
        mutation.mutate(params)
    }
    const queryClient = useQueryClient()
    const deleteMutation = useMutation(deleteMyBlogPost, {
        onSuccess: () => {
            queryClient.invalidateQueries('blogposts')
        }
    })

    const createBlogPostMutation = useMutation(createBlogPost, {
        onSuccess: (data) => {
            console.log("data", data);
            const newPostId = data.id;
            queryClient.invalidateQueries('blogposts');
            router.push(`/portal/my_posts/draft?id=${newPostId}`);
        }
    });

    const handleCreateClick = () => {
        createBlogPostMutation.mutate({
            authToken: cookies.user_token,
        });
    };


    const handleDelete = (id: number) => {
        setConfirmationModalData({
            onMouseDownFunc: () => mutationFunc(deleteMutation, { id: id.toString(), authToken: cookies.user_token }),
            warningText: "Deleting a blog post means it will be completely removed forever. This action cannot be undone.",
        });
    };

    const handleEdit = (id: number) => {
        router.push(`/portal/my_posts/draft?id=${id}`);
    };

    const blogPostsRows = rows.map((blogPost) => {
        const { id, title, author, status, created_at } = blogPost;
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
                        {' '}Go to blog
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
                    emptyStateMsg={<BlogEmptyState />}
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
