"use client";

import { useMutation, useQuery, useQueryClient } from "react-query";
import { ChangeEvent, useState, useEffect, Suspense } from "react";
import { getJobPost, updateMyJobPost, deleteJobPost } from "../../../queries";
import { JobPost } from "../../../types";
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
    const [JobPostTitleString, setJobPostTitleString] = useState<string>("");
    const [JobPostContentString, setJobPostContentString] = useState<string>("");
    const [errorText, setErrorText] = useState<string>("");
    const [submitted, setSubmitted] = useState<boolean>(false);
    const [showModal, setShowModal] = useState<boolean>(false);
    const [renderedContent, setRenderedContent] = useState<string>("");
    const [showSuccessIcon, setShowSuccessIcon] = useState<boolean>(false);
    const queryClient = useQueryClient();

    const { data: jobPostData, isLoading, isError } = useQuery<JobPost, Error>({
        queryKey: ['jobpost', postId],
        queryFn: () => getJobPost({ authToken: cookies.user_token, id: postId as string }),
        enabled: !!postId,
        retry: (failureCount, error): boolean => {
            if (error.message.includes("401")) {
                return false;
            }
            return true;
        },
    });

    useEffect(() => {
        if (jobPostData) {
            setJobPostTitleString(jobPostData?.title || "");
            setJobPostContentString(jobPostData?.content || "");
        }
    }, [jobPostData]);

    const saveJobPostMutation = useMutation(updateMyJobPost, {
        onSuccess: () => {
            setErrorText("");
            queryClient.invalidateQueries('jobposts');
            queryClient.invalidateQueries('jobpost');

            setShowSuccessIcon(true);
            setTimeout(() => setShowSuccessIcon(false), 1000);
        },
        onError: (e: Error) => {
            setErrorText(e.message);
        }
    });

    const publishJobPostMutation = useMutation(updateMyJobPost, {
        onSuccess: () => {
            setErrorText("");
            setSubmitted(true);
            queryClient.invalidateQueries('jobposts');
            queryClient.invalidateQueries('jobpost');
        },
        onError: (e: Error) => {
            setErrorText(e.message);
        }
    });

    const deleteJobPostMutation = useMutation(deleteJobPost, {
        onSuccess: () => {
            setErrorText("");
            queryClient.invalidateQueries('jobposts');
            router.push("/portal/my_posts");
        },
        onError: (e: Error) => {
            setErrorText(e.message);
        }
    });

    const handleJobPostTitleChange = (event: ChangeEvent<HTMLInputElement>) => {
        setJobPostTitleString(event.target.value);
    };

    const handleJobPostContentChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
        setJobPostContentString(event.target.value);
    };

    const handleShowClick = async () => {
        const processedContent = await remark().use(html).process(JobPostContentString);
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
        return <div>No job post ID found.</div>;
    }

    const handleDiscardClick = () => {
        deleteJobPostMutation.mutate({
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
                    title="Edit Job Post"
                    controls={
                        <>
                            <Button
                                appearance="base"
                                disabled={JobPostTitleString === "" || JobPostContentString === ""}
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
                            value={JobPostTitleString}
                            onChange={handleJobPostTitleChange}
                            style={{
                                backgroundColor: "white",
                                border: "1px solid #ddd",
                                padding: "10px"
                            }}
                        />
                        <Textarea
                            placeholder={`Write your job post here. You can use markdown to format your text.`}
                            rows={20}
                            value={JobPostContentString} // Use the fetched content
                            onChange={handleJobPostContentChange}
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
                            disabled={JobPostTitleString === "" || JobPostContentString === ""}
                            onClick={(event) => {
                                event.preventDefault();
                                publishJobPostMutation.mutate({
                                    authToken: cookies.user_token,
                                    id: postId,
                                    status: "published",
                                    title: JobPostTitleString,
                                    content: JobPostContentString
                                });
                            }}
                        >
                            Publish
                        </Button>
                        <Button
                            appearance="base"
                            name="save"
                            disabled={JobPostTitleString === "" || JobPostContentString === ""}
                            onClick={(event) => {
                                event.preventDefault();
                                saveJobPostMutation.mutate({
                                    authToken: cookies.user_token,
                                    id: postId,
                                    status: "draft",
                                    title: JobPostTitleString,
                                    content: JobPostContentString
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
                        <h2>{JobPostTitleString}</h2>
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
