const okareo_sdk = require('okareo-ts-sdk');

const okareo = new okareo_sdk.Okareo({api_key:process.env.OKAREO_API_KEY, endpoint: process.env.OKAREO_BASE_URL});

const main = async () => {
    try {
        const projects = await okareo.getProjects();
        console.log("javascript", projects);
    } catch (error) {
        console.error(error);
    }
}

main();