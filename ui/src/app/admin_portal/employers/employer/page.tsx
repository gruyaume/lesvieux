"use client";

import { useSearchParams, useRouter } from "next/navigation";
import { Suspense, useState, useEffect } from "react";
import { numEmployerAccounts, getEmployer } from "../../../queries";
import { useAuth } from "../../auth/authContext";
import { useCookies } from "react-cookie";
import { Panel } from "@canonical/react-components";

function EmployerContent() {
    const auth = useAuth();
    const [cookies] = useCookies(['user_token']);
    const searchParams = useSearchParams();
    const employer_id = searchParams.get("employer_id");
    const router = useRouter();
    const [accountCount, setAccountCount] = useState<number | null>(null);
    const [employerName, setEmployerName] = useState<string | null>(null);

    useEffect(() => {
        if (employer_id && cookies.user_token) {
            numEmployerAccounts({ authToken: cookies.user_token, employerId: employer_id })
                .then((count) => {
                    setAccountCount(count);
                })
                .catch((error) => {
                    console.error("Failed to fetch employer accounts:", error);
                    setAccountCount(null);
                });
        }
    }, [employer_id, cookies.user_token]);

    useEffect(() => {
        if (employer_id && cookies.user_token) {
            getEmployer({ authToken: cookies.user_token, id: employer_id })
                .then((employer) => {
                    setEmployerName(employer.name);
                })
                .catch((error) => {
                    console.error("Failed to fetch employer name:", error);
                    setEmployerName(null);
                });
        }
    }, [employer_id, cookies.user_token]);


    if (!employer_id) {
        return <div>Employer not found.</div>;
    }

    if (!auth.user) {
        return null;
    }

    const handleBoxClick = () => {
        router.push(`/admin_portal/employers/employer/accounts?employer_id=${employer_id}`);
    };

    return (
        <Panel
            title={`Employers / ${employerName}`}
            stickyHeader
        >
            <div style={{ padding: "20px" }}>
                <div style={{ marginTop: "20px", display: "flex", flexDirection: "column", alignItems: "flex-start" }}>
                    <p style={{ fontSize: "18px", fontWeight: "bold", margin: "0 0 8px 0", color: "#555" }}>
                        Accounts
                    </p>
                    <div
                        onClick={handleBoxClick}
                        style={{
                            display: "flex",
                            justifyContent: "center",
                            alignItems: "center",
                            width: "150px",
                            height: "150px",
                            border: "2px solid #0078d4",
                            borderRadius: "12px",
                            backgroundColor: "#f4f4f4",
                            cursor: "pointer", // Indicates clickable element
                        }}
                    >
                        {accountCount !== null ? (
                            <p style={{ fontSize: "48px", fontWeight: "bold", margin: 0, color: "#333" }}>
                                {accountCount}
                            </p>
                        ) : (
                            <p style={{ fontSize: "18px", margin: 0, color: "#666" }}>
                                Loading...
                            </p>
                        )}
                    </div>
                </div>
            </div>
        </Panel>
    );
}

export default function Employer() {
    return (
        <Suspense fallback={<div>Loading...</div>}>
            <EmployerContent />
        </Suspense>
    );
}
