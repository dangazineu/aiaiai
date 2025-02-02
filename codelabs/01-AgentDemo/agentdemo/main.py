from datetime import datetime

from openai import OpenAI
from pydantic.v1 import BaseModel

from agent import OpenAIAgent
from tool import Tool

from langchain_core.utils.function_calling import convert_to_openai_tool

# === TOOL DEFINITIONS ===

class Expense(BaseModel):
    description: str
    net_amount: float
    gross_amount: float
    tax_rate: float
    date: datetime


def add_expense_func(**kwargs):
    return f"Added expense: {kwargs} to the database."

class ReportTool(BaseModel):
    report: str = None


def report_func(report: str = None):
    return f"Reported: {report}"

class DateTool(BaseModel):
    x: str = None

def main():
    add_expense_tool = Tool(
        name="add_expense_tool",
        model=Expense,
        function=add_expense_func
    )
    report_tool = Tool(
        name="report_tool",
        model=ReportTool,
        function=report_func
    )
    get_date_tool = Tool(
        name="get_current_date",
        model=DateTool,
        function=lambda: datetime.now().strftime("%Y-%m-%d"),
        validate_missing=False
    )

    tools = [add_expense_tool, report_tool, get_date_tool]
    client = OpenAI()
    model_name = "gpt-3.5-turbo-0125"
    agent = OpenAIAgent(tools, client, model_name=model_name, verbose=True)

    user_input = "I have spend 5$ on a coffee today please track my expense. The tax rate is 0.2"

    agent.run(user_input)


if __name__ == "__main__":
    main()