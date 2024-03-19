import { Okareo } from 'okareo-ts-sdk';

const okareo = new Okareo({api_key:process.env.OKAREO_API_KEY, endpoint: process.env.OKAREO_BASE_URL});

const main = async () => {
    try {
        const projects = await okareo.getProjects();
        console.log("Typescript", projects);
    } catch (error) {
        console.error(error);
    }
}

main();