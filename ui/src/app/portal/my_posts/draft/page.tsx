"use client";

import { useMutation, useQuery, useQueryClient } from "react-query";
import { ChangeEvent, useState, useEffect, Suspense } from "react";
import { getBlogPost, updateMyBlogPost, deleteBlogPost } from "../../../queries";
import { BlogPost } from "../../../types";
import { useCookies } from "react-cookie";
import { useRouter, useSearchParams } from "next/navigation";
import { Input, Button, Form, Panel, Textarea, Icon, Modal } from "@canonical/react-components";
import { useAuth } from "../../auth/authContext";
import { remark } from 'remark';
import html from 'remark-html';

function DraftContent() {
    const searchParams = useSearchParams();
    const postId = searchParams.get("id");
    const router = useRouter();
    const auth = useAuth();
    const [cookies] = useCookies(['user_token']);
    const [BlogPostTitleString, setBlogPostTitleString] = useState<string>("");
    const [BlogPostContentString, setBlogPostContentString] = useState<string>("");
    const [errorText, setErrorText] = useState<string>("");
    const [submitted, setSubmitted] = useState<boolean>(false);
    const [showModal, setShowModal] = useState<boolean>(false);
    const [renderedContent, setRenderedContent] = useState<string>("");
    const [showSuccessIcon, setShowSuccessIcon] = useState<boolean>(false);
    const queryClient = useQueryClient();

    const { data: blogPostData, isLoading, isError } = useQuery<BlogPost, Error>({
        queryKey: ['blogpost', postId],
        queryFn: () => getBlogPost({ authToken: cookies.user_token, id: postId as string }),
        enabled: !!postId,
        retry: (failureCount, error): boolean => {
            if (error.message.includes("401")) {
                return false;
            }
            return true;
        },
    });

    useEffect(() => {
        if (blogPostData) {
            setBlogPostTitleString(blogPostData?.title || "");
            setBlogPostContentString(blogPostData?.content || "");
        }
    }, [blogPostData]);

    const saveBlogPostMutation = useMutation(updateMyBlogPost, {
        onSuccess: () => {
            setErrorText("");
            queryClient.invalidateQueries('blogposts');
            queryClient.invalidateQueries('blogpost');

            setShowSuccessIcon(true);
            setTimeout(() => setShowSuccessIcon(false), 1000);
        },
        onError: (e: Error) => {
            setErrorText(e.message);
        }
    });

    const publishBlogPostMutation = useMutation(updateMyBlogPost, {
        onSuccess: () => {
            setErrorText("");
            setSubmitted(true);
            queryClient.invalidateQueries('blogposts');
            queryClient.invalidateQueries('blogpost');
        },
        onError: (e: Error) => {
            setErrorText(e.message);
        }
    });

    const deleteBlogPostMutation = useMutation(deleteBlogPost, {
        onSuccess: () => {
            setErrorText("");
            queryClient.invalidateQueries('blogposts');
            router.push("/portal/my_posts");
        },
        onError: (e: Error) => {
            setErrorText(e.message);
        }
    });

    const handleBlogPostTitleChange = (event: ChangeEvent<HTMLInputElement>) => {
        setBlogPostTitleString(event.target.value);
    };

    const handleBlogPostContentChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
        setBlogPostContentString(event.target.value);
    };

    const handleShowClick = async () => {
        const processedContent = await remark().use(html).process(BlogPostContentString);
        setRenderedContent(processedContent.toString());
        setShowModal(true);
    };

    useEffect(() => {
        if (submitted) {
            router.push("/portal/my_posts");
        }
    }, [submitted, router]);

    if (!postId) {
        // Handle the case where postId is not available
        return <div>No blog post ID found.</div>;
    }

    const handleDiscardClick = () => {
        deleteBlogPostMutation.mutate({
            authToken: cookies.user_token,
            id: postId
        });
    };

    const today = new Date().toISOString().split('T')[0]; // Format YYYY-MM-DD

    if (!auth.user) {
        return null;
    }

    return (
        <div style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            minHeight: "100vh",
            padding: "5vh 0",
            backgroundColor: "#f5f5f5",
            boxSizing: "border-box"
        }}>
            <div style={{
                width: "90%",
                maxWidth: "950px",
                maxHeight: "85vh",
                backgroundColor: "white",
                boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
                borderRadius: "8px",
                display: "flex",
                flexDirection: "column",
                justifyContent: "space-between",
                padding: "20px",
                boxSizing: "border-box"
            }}>
                <Panel
                    title="Edit Blog Post"
                    controls={
                        <>
                            <Button
                                appearance="base"
                                disabled={BlogPostTitleString === "" || BlogPostContentString === ""}
                                onClick={handleShowClick}>
                                <Icon name="show" />
                            </Button>
                            <Button
                                appearance="base"
                                onClick={handleDiscardClick}>
                                <Icon name="delete" />
                            </Button>
                            <Button
                                appearance="base"
                                onClick={() => { router.push("/portal/my_posts") }}>
                                <Icon name="external-link" />
                            </Button>
                        </>
                    }
                >
                    <Form>
                        <Input
                            id="InputTitle"
                            type="text"
                            placeholder="Title"
                            value={BlogPostTitleString}
                            onChange={handleBlogPostTitleChange}
                            style={{
                                backgroundColor: "white",
                                border: "1px solid #ddd",
                                padding: "10px"
                            }}
                        />
                        <Textarea
                            placeholder={`Write your blog post here. You can use markdown to format your text.`}
                            rows={20}
                            value={BlogPostContentString} // Use the fetched content
                            onChange={handleBlogPostContentChange}
                            error={errorText}
                            style={{
                                backgroundColor: "white",
                                border: "1px solid #ddd",
                                padding: "10px"
                            }}
                        />
                        <Button
                            appearance="positive"
                            name="submit"
                            disabled={BlogPostTitleString === "" || BlogPostContentString === ""}
                            onClick={(event) => {
                                event.preventDefault();
                                publishBlogPostMutation.mutate({
                                    authToken: cookies.user_token,
                                    id: postId,
                                    status: "published",
                                    title: BlogPostTitleString,
                                    content: BlogPostContentString
                                });
                            }}
                        >
                            Publish
                        </Button>
                        <Button
                            appearance="base"
                            name="save"
                            disabled={BlogPostTitleString === "" || BlogPostContentString === ""}
                            onClick={(event) => {
                                event.preventDefault();
                                saveBlogPostMutation.mutate({
                                    authToken: cookies.user_token,
                                    id: postId,
                                    status: "draft",
                                    title: BlogPostTitleString,
                                    content: BlogPostContentString
                                });
                            }}
                        >
                            Save draft {"  "}
                            {showSuccessIcon && <Icon name="success" />}
                        </Button>
                    </Form>
                </Panel>
                {showModal && (
                    <Modal
                        title="Preview"
                        close={() => setShowModal(false)}
                    >
                        <h2>{BlogPostTitleString}</h2>
                        <h5>By: {auth.user.username}</h5>
                        <h6>{today}</h6>
                        <div
                            dangerouslySetInnerHTML={{ __html: renderedContent }}
                        />
                    </Modal>
                )}
            </div>
        </div>
    );
}

export default function Draft() {
    return (
        <Suspense fallback={<div>Loading...</div>}>
            <DraftContent />
        </Suspense>
    );
}
