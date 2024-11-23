"use client"

import { Application, AppMain, Navigation, Panel } from "@canonical/react-components";
import Logo from "./components/logo";

export default function ReaderNavigation({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <Application>
            <AppMain>
                <Navigation
                    items={[
                    ]}
                    itemsRight={[
                        {
                            alignRight: true,
                            items: [
                                {
                                    label: 'Applicant',
                                    url: '#'
                                },
                                {
                                    label: 'Employer',
                                    url: '/employer_portal/login'
                                },
                                {
                                    label: 'Admin',
                                    url: '/admin_portal/login'
                                }
                            ],
                            label: 'Portal'
                        }
                    ]}
                    logo={<Logo />} />
                <div style={{
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    padding: "20px",
                }}>
                    <div style={{
                        width: "100%",
                        maxWidth: "800px",
                        minWidth: "600px",
                        padding: "20px",
                    }}>
                        <Panel>
                            <p>
                                {children}
                            </p>
                        </Panel>
                    </div>
                </div>
            </AppMain>
        </Application>
    )
}
