#!/usr/bin/env python3

import os
import uuid
import random
import string
from okareo import Okareo
from okareo_api_client.models import ScenarioSetCreate, ScenarioSetResponse, SeedData, ScenarioType
from okareo.model_under_test import OpenAIModel
from okareo_api_client.models.test_run_type import TestRunType

OKAREO_API_KEY = os.environ["OKAREO_API_KEY"]
OPENAI_API_KEY = os.environ["OPENAI_API_KEY"]
OKAREO_RUN_ID = os.environ["OKAREO_RUN_ID"]

okareo = Okareo(OKAREO_API_KEY)

USER_PROMPT_TEMPLATE = "{input}"

SUMMARIZATION_CONTEXT_TEMPLATE = """
You will be provided with text.
Summarize the text in 1-3 simple sentences.
If the text's title and the author's name are available, refer to both in the summary.
Your goal is to make the summary clear to a middle school student.
"""

# Text for the summarization model to process
text_1 = """
I Wandered Lonely as a Cloud

I wandered lonely as a cloud
That floats on high o'er vales and hills,
When all at once I saw a crowd,
A host, of golden daffodils;
Beside the lake, beneath the trees,
Fluttering and dancing in the breeze.

Continuous as the stars that shine
And twinkle on the milky way,
They stretched in never-ending line
Along the margin of a bay:
Ten thousand saw I at a glance,
Tossing their heads in sprightly dance.

The waves beside them danced; but they
Out-did the sparkling waves in glee:
A poet could not but be gay,
In such a jocund company:
I gazed-and gazed-but little thought
What wealth the show to me had brought:

For oft, when on my couch I lie
In vacant or in pensive mood,
They flash upon that inward eye
Which is the bliss of solitude;
And then my heart with pleasure fills,
And dances with the daffodils.

-- by William Wordsworth
"""

text_2 = """
All Things Are Current Found

ALL things are current found
On earthly ground,
Spirits and elements
Have their descents.

Night and day, year on year,
High and low, far and near,
These are our own aspects,
These are our own regrets.

Ye gods of the shore,
Who abide evermore,
I see you far headland,
Stretching on either hand;

I hear the sweet evening sounds
From your undecaying grounds;
Cheat me no more with time,
Take me to your clime.

-- by Henry D. Thoreau
"""

text_3 = """
Hope Is the Thing with Feathers

"Hope" is the thing with feathers-
That perches in the soul-
And sings the tune without the words-
And never stops-at all-

And sweetest-in the Gale-is heard-
And sore must be the storm-
That could abash the little Bird
That kept so many warm-

I've heard it in the chillest land-
And on the strangest Sea-
Yet, never, in Extremity,
It asked a crumb-of Me.

-- by Emily Dickinson
"""

scenario_set_create = ScenarioSetCreate(
    name=f"{OKAREO_RUN_ID} - Scenario",
    number_examples=1,
    generation_type=ScenarioType.SEED,
    seed_data=[
        SeedData(
            input_=text_1,  
            result=str(uuid.uuid4())
        ),
        SeedData(
            input_=text_2,  
            result=str(uuid.uuid4())
        ),
        SeedData(
            input_=text_3, 
            result=str(uuid.uuid4())
        )
    ],
)
scenario = okareo.create_scenario_set(scenario_set_create)

print(scenario.additional_properties['app_link'])

mut_name = f"{OKAREO_RUN_ID} - MUT"
eval_name = f"{OKAREO_RUN_ID} - EVAL"

model_under_test = okareo.register_model(
    name=mut_name,
    tags=[OKAREO_RUN_ID],
    model=OpenAIModel(
        model_id="gpt-3.5-turbo",
        temperature=0,
        system_prompt_template=SUMMARIZATION_CONTEXT_TEMPLATE,
        user_prompt_template=USER_PROMPT_TEMPLATE,
    ),
)

evaluation = model_under_test.run_test(
    name=eval_name,
    scenario=scenario,
    api_key=OPENAI_API_KEY,
    test_run_type=TestRunType.NL_GENERATION,
    calculate_metrics=True,
)

print(evaluation.additional_properties['app_link'])