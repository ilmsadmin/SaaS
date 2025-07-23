-- LMS Service Database Migration
-- Create LMS tables with multi-tenant support

-- Categories table for course categorization
CREATE TABLE IF NOT EXISTS course_categories (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    color VARCHAR(7) DEFAULT '#007bff',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Courses table for course management
CREATE TABLE IF NOT EXISTS courses (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    category_id INTEGER REFERENCES course_categories(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    short_description VARCHAR(500),
    thumbnail_url TEXT,
    video_url TEXT,
    level VARCHAR(20) DEFAULT 'beginner', -- beginner, intermediate, advanced
    status VARCHAR(20) DEFAULT 'draft', -- draft, published, archived
    duration_hours INTEGER DEFAULT 0,
    max_students INTEGER,
    price DECIMAL(10,2) DEFAULT 0.00,
    is_free BOOLEAN DEFAULT true,
    is_featured BOOLEAN DEFAULT false,
    requirements TEXT,
    what_you_learn TEXT,
    instructor_id VARCHAR(255),
    instructor_name VARCHAR(255),
    language VARCHAR(10) DEFAULT 'en',
    tags TEXT,
    enrolled_count INTEGER DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.00,
    rating_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Course sections/modules
CREATE TABLE IF NOT EXISTS course_sections (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Course lessons/lectures
CREATE TABLE IF NOT EXISTS course_lessons (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
    section_id INTEGER REFERENCES course_sections(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content_type VARCHAR(20) DEFAULT 'video', -- video, text, quiz, assignment
    content_url TEXT,
    content_text TEXT,
    duration_minutes INTEGER DEFAULT 0,
    sort_order INTEGER DEFAULT 0,
    is_preview BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Student enrollments
CREATE TABLE IF NOT EXISTS enrollments (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
    student_id VARCHAR(255) NOT NULL,
    student_name VARCHAR(255),
    student_email VARCHAR(255),
    enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completion_date TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active', -- active, completed, dropped, suspended
    progress_percentage DECIMAL(5,2) DEFAULT 0.00,
    last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    certificate_issued BOOLEAN DEFAULT false,
    certificate_url TEXT,
    payment_status VARCHAR(20) DEFAULT 'pending', -- pending, paid, refunded
    payment_amount DECIMAL(10,2) DEFAULT 0.00,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Student lesson progress
CREATE TABLE IF NOT EXISTS lesson_progress (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    enrollment_id INTEGER REFERENCES enrollments(id) ON DELETE CASCADE,
    lesson_id INTEGER REFERENCES course_lessons(id) ON DELETE CASCADE,
    student_id VARCHAR(255) NOT NULL,
    status VARCHAR(20) DEFAULT 'not_started', -- not_started, in_progress, completed
    progress_percentage DECIMAL(5,2) DEFAULT 0.00,
    time_spent_minutes INTEGER DEFAULT 0,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quizzes and assessments
CREATE TABLE IF NOT EXISTS quizzes (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
    lesson_id INTEGER REFERENCES course_lessons(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    instructions TEXT,
    time_limit_minutes INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 0,
    passing_score DECIMAL(5,2) DEFAULT 70.00,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quiz questions
CREATE TABLE IF NOT EXISTS quiz_questions (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    quiz_id INTEGER REFERENCES quizzes(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    question_type VARCHAR(20) DEFAULT 'multiple_choice', -- multiple_choice, true_false, short_answer, essay
    options JSON, -- Store answer options for multiple choice
    correct_answer TEXT,
    points DECIMAL(5,2) DEFAULT 1.00,
    explanation TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Student quiz attempts
CREATE TABLE IF NOT EXISTS quiz_attempts (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    quiz_id INTEGER REFERENCES quizzes(id) ON DELETE CASCADE,
    student_id VARCHAR(255) NOT NULL,
    enrollment_id INTEGER REFERENCES enrollments(id) ON DELETE CASCADE,
    attempt_number INTEGER DEFAULT 1,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    time_spent_minutes INTEGER DEFAULT 0,
    score DECIMAL(5,2) DEFAULT 0.00,
    max_score DECIMAL(5,2) DEFAULT 0.00,
    percentage DECIMAL(5,2) DEFAULT 0.00,
    status VARCHAR(20) DEFAULT 'in_progress', -- in_progress, completed, abandoned
    answers JSON, -- Store student answers
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Course reviews and ratings
CREATE TABLE IF NOT EXISTS course_reviews (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
    student_id VARCHAR(255) NOT NULL,
    enrollment_id INTEGER REFERENCES enrollments(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    review_text TEXT,
    is_featured BOOLEAN DEFAULT false,
    is_approved BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Assignments and submissions
CREATE TABLE IF NOT EXISTS assignments (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
    lesson_id INTEGER REFERENCES course_lessons(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    instructions TEXT,
    due_date TIMESTAMP,
    max_points DECIMAL(5,2) DEFAULT 100.00,
    submission_type VARCHAR(20) DEFAULT 'file', -- file, text, url
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Assignment submissions
CREATE TABLE IF NOT EXISTS assignment_submissions (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    assignment_id INTEGER REFERENCES assignments(id) ON DELETE CASCADE,
    student_id VARCHAR(255) NOT NULL,
    enrollment_id INTEGER REFERENCES enrollments(id) ON DELETE CASCADE,
    submission_text TEXT,
    submission_url TEXT,
    file_path TEXT,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    graded_at TIMESTAMP,
    grade DECIMAL(5,2),
    feedback TEXT,
    status VARCHAR(20) DEFAULT 'submitted', -- submitted, graded, returned
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_course_categories_tenant_id ON course_categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_courses_tenant_id ON courses(tenant_id);
CREATE INDEX IF NOT EXISTS idx_courses_category_id ON courses(category_id);
CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);
CREATE INDEX IF NOT EXISTS idx_courses_instructor_id ON courses(instructor_id);
CREATE INDEX IF NOT EXISTS idx_course_sections_tenant_id ON course_sections(tenant_id);
CREATE INDEX IF NOT EXISTS idx_course_sections_course_id ON course_sections(course_id);
CREATE INDEX IF NOT EXISTS idx_course_lessons_tenant_id ON course_lessons(tenant_id);
CREATE INDEX IF NOT EXISTS idx_course_lessons_course_id ON course_lessons(course_id);
CREATE INDEX IF NOT EXISTS idx_course_lessons_section_id ON course_lessons(section_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_tenant_id ON enrollments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_course_id ON enrollments(course_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_student_id ON enrollments(student_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_status ON enrollments(status);
CREATE INDEX IF NOT EXISTS idx_lesson_progress_tenant_id ON lesson_progress(tenant_id);
CREATE INDEX IF NOT EXISTS idx_lesson_progress_enrollment_id ON lesson_progress(enrollment_id);
CREATE INDEX IF NOT EXISTS idx_lesson_progress_lesson_id ON lesson_progress(lesson_id);
CREATE INDEX IF NOT EXISTS idx_lesson_progress_student_id ON lesson_progress(student_id);
CREATE INDEX IF NOT EXISTS idx_quizzes_tenant_id ON quizzes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_quizzes_course_id ON quizzes(course_id);
CREATE INDEX IF NOT EXISTS idx_quiz_questions_tenant_id ON quiz_questions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_quiz_questions_quiz_id ON quiz_questions(quiz_id);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_tenant_id ON quiz_attempts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_quiz_id ON quiz_attempts(quiz_id);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_student_id ON quiz_attempts(student_id);
CREATE INDEX IF NOT EXISTS idx_course_reviews_tenant_id ON course_reviews(tenant_id);
CREATE INDEX IF NOT EXISTS idx_course_reviews_course_id ON course_reviews(course_id);
CREATE INDEX IF NOT EXISTS idx_assignments_tenant_id ON assignments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_assignments_course_id ON assignments(course_id);
CREATE INDEX IF NOT EXISTS idx_assignment_submissions_tenant_id ON assignment_submissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_assignment_submissions_assignment_id ON assignment_submissions(assignment_id);
CREATE INDEX IF NOT EXISTS idx_assignment_submissions_student_id ON assignment_submissions(student_id);

-- Create triggers for updating timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_course_categories_updated_at BEFORE UPDATE ON course_categories FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_courses_updated_at BEFORE UPDATE ON courses FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_course_sections_updated_at BEFORE UPDATE ON course_sections FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_course_lessons_updated_at BEFORE UPDATE ON course_lessons FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_enrollments_updated_at BEFORE UPDATE ON enrollments FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_lesson_progress_updated_at BEFORE UPDATE ON lesson_progress FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_quizzes_updated_at BEFORE UPDATE ON quizzes FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_quiz_questions_updated_at BEFORE UPDATE ON quiz_questions FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_quiz_attempts_updated_at BEFORE UPDATE ON quiz_attempts FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_course_reviews_updated_at BEFORE UPDATE ON course_reviews FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_assignments_updated_at BEFORE UPDATE ON assignments FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_assignment_submissions_updated_at BEFORE UPDATE ON assignment_submissions FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
