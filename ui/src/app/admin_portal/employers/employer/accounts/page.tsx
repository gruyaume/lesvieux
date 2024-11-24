"use client";

import { useSearchParams } from "next/navigation";
import { Suspense, useState, useEffect } from "react";
import { useQuery } from "react-query";
import { useAuth } from "../../../auth/authContext";
import { useRouter } from "next/navigation";
import { useCookies } from "react-cookie";
import { EmployerUsersTable } from "./table";
import { listEmployerAccounts, getEmployer } from "../../../../queries";
import { UserEntry } from "../../../../types";

function EmployerAccountsContent() {
    const auth = useAuth();
    const [cookies, , removeCookie] = useCookies(["user_token"]);
    const searchParams = useSearchParams();
    const employer_id = searchParams.get("employer_id");
    const router = useRouter();
    const [employerName, setEmployerName] = useState<string | null>(null);

    // Redirect if user is not logged in
    if (!cookies.user_token) {
        router.push("/admin_portal/login");
    }

    // Fetch employer users
    const employerUsersQuery = useQuery<UserEntry[], Error>({
        queryKey: ["employer_users", cookies.user_token, employer_id],
        queryFn: () => {
            if (!employer_id) {
                return Promise.reject(new Error("Employer ID is missing"));
            }
            return listEmployerAccounts({ authToken: cookies.user_token, employerId: employer_id });
        },
        retry: (failureCount, error) => {
            if (error.message.includes("401")) {
                return false;
            }
            return true;
        },
    });

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

    if (employerUsersQuery.isLoading) {
        return <div>Loading...</div>;
    }

    if (employerUsersQuery.isError) {
        if (employerUsersQuery.error.message.includes("401")) {
            removeCookie("user_token");
            router.push("/admin_portal/login");
            return null;
        }
        return <div>Error: {employerUsersQuery.error.message}</div>;
    }

    if (employerName === null) {
        return <div>Employer not found.</div>;
    }

    const employerUsers = employerUsersQuery.data || [];
    return <EmployerUsersTable employerUsers={employerUsers} employerId={employer_id} employerName={employerName} />;
}

export default function EmployerAccounts() {
    return (
        <Suspense fallback={<div>Loading...</div>}>
            <EmployerAccountsContent />
        </Suspense>
    );
}
