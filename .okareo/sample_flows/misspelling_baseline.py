from basic_generator_run import run_basic_generator
from okareo_api_client.models import ScenarioType

run_basic_generator(run_name='MISSPELLING',
                    scenario_type=ScenarioType.COMMON_MISSPELLINGS)
