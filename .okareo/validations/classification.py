#!/usr/bin/env python3

import os
from okareo import Okareo
from okareo_api_client.models import ScenarioSetCreate, SeedData, ScenarioType
from okareo.model_under_test import OpenAIModel
from okareo_api_client.models.test_run_type import TestRunType

import random
import string

def generate_random_string(length):
    alphabet = string.ascii_letters + string.digits
    return ''.join(random.choice(alphabet) for _ in range(length))

random_string = generate_random_string(6)

OKAREO_API_KEY = os.environ["OKAREO_API_KEY"]
OPENAI_API_KEY = os.environ["OPENAI_API_KEY"]
OKAREO_RUN_ID = os.environ["OKAREO_RUN_ID"]

okareo = Okareo(OKAREO_API_KEY)

# Simple adhoc classifier using OpenAI'a GPT 3.5 Turbo model

USER_PROMPT_TEMPLATE = "{input}"

CLASSIFICATION_CONTEXT_TEMPLATE = """
You will be provided a question from a customer.
Classify the question into a customer category and sub-category.
Provide the output with only the category name.

Categories: Technical Support, Billing, Account Management, General Inquiry, Unknown

Sub-Categories for Technical Support:
Troubleshooting
Product features
Product updates
Integration options
Found a problem

Sub-Categories for Billing:
Unsubscribe
Upgrade
Explain my bill
Change payment
Dispute a charge

Sub-Categories for Account Management:
Add a team member
Change or Update details
Password reset
Close account
Security

Sub-Categories for General Inquiry:
Contact sales
Product information
Pricing
Feedback
Speak to a human
"""

scenario_set_create = ScenarioSetCreate(
    name=f"Github Action Test ID {OKAREO_RUN_ID} - {random_string} - Scenario",
    number_examples=1,
    generation_type=ScenarioType.SEED,
    seed_data=[
        SeedData(
            input_="Can I connect to my SalesForce?",  
            result="Technical Support"
        ),
        SeedData(
            input_="Do you have a way to send marketing emails?",  
            result="Technical Support"
        ),
        SeedData(
            input_="Can I get invoiced instead of using a credit card?", 
            result="Billing"
        ),
        SeedData(
            input_="My CRM integration is not working.", 
            result="Technical Support"
        ),
        SeedData(
            input_="Do you have SOC II tpye 2 certification?", 
            result="Account Management"
        ),
        SeedData(
            input_="I like the product.  Please connect me to your enterprise team.", 
            result="General Inquiry"
        )
    ],
)
scenario = okareo.create_scenario_set(scenario_set_create)
print('Scenario Link: ', scenario.app_link)

# Establish the model that is being evaluated, at minimum this is a named model for future reference
model_under_test = okareo.register_model(
    name=f"Github Action Test ID {OKAREO_RUN_ID} - {random_string} - MUT",
    tags=[OKAREO_RUN_ID],
    model=OpenAIModel(
        model_id="gpt-3.5-turbo",
        temperature=0,
        system_prompt_template=CLASSIFICATION_CONTEXT_TEMPLATE,
        user_prompt_template=USER_PROMPT_TEMPLATE,
    ),
)

# run the test and call the model for each item in the scenario set
evaluation = model_under_test.run_test(
    name=f"Github Action Test ID {OKAREO_RUN_ID} - {random_string} - EVAL",
    scenario=scenario,
    api_key=OPENAI_API_KEY,
    test_run_type=TestRunType.MULTI_CLASS_CLASSIFICATION,
    calculate_metrics=True,
)

print('Evaluation Link: ', evaluation.app_link)