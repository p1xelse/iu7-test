import pytest
from sqlalchemy import create_engine
from sqlalchemy.orm import Session
from sqlalchemy import text

@pytest.fixture(scope='module')
def api_url():
    return 'http://localhost:8080'

@pytest.fixture(autouse=True)
def db_all_tables():
    # setup
    engine = create_engine(
        'postgresql://test:test@localhost:13081/postgres', 
        isolation_level="READ UNCOMMITTED")
    connection = engine.connect()
    session = Session(bind=connection)
    session.execute(text("truncate table users cascade;"))
    session.execute(text("truncate table project cascade;"))
    session.execute(text("truncate table tag cascade;"))
    session.execute(text("truncate table friend_relation cascade;"))
    session.execute(text("truncate table tag_entry cascade;"))
    session.execute(text("truncate table goal cascade;"))
    session.execute(text("truncate table entry cascade;"))
    session.commit()

    # return session
    yield session

    # tear down
    session.execute(text("truncate table users cascade;"))
    session.execute(text("truncate table project cascade;"))
    session.execute(text("truncate table tag cascade;"))
    session.execute(text("truncate table friend_relation cascade;"))
    session.execute(text("truncate table tag_entry cascade;"))
    session.execute(text("truncate table goal cascade;"))
    session.execute(text("truncate table entry cascade;"))
    session.commit()