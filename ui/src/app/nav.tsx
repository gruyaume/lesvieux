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
                <Navigation items={[]} logo={<Logo />} />
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
