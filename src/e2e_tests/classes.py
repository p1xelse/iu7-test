from sqlalchemy import Column, Integer, String, Text, ForeignKey, DateTime, Float, Boolean
from sqlalchemy.orm import relationship
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()

class User(Base):
    __tablename__ = 'users'
    
    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String(35), nullable=False)
    email = Column(String(254), nullable=False, unique=True)
    about = Column(Text, default='')
    role = Column(String(20), default='user')
    password = Column(String(128), nullable=False)

class Entry(Base):
    __tablename__ = 'entry'

    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, ForeignKey('users.id'), nullable=False)
    project_id = Column(Integer, ForeignKey('project.id', ondelete='CASCADE'))
    description = Column(Text, default='')
    time_start = Column(DateTime, nullable=False)
    time_end = Column(DateTime, nullable=False)

    user = relationship('User', backref='entries')
    project = relationship('Project')

class Project(Base):
    __tablename__ = 'project'

    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, ForeignKey('users.id'), nullable=False)
    name = Column(String(35), nullable=False)
    about = Column(Text, default='')
    color = Column(String(10), nullable=False)
    is_private = Column(Boolean, nullable=False)
    total_count_hours = Column(Float, default=0)

    user = relationship('User', backref='projects')