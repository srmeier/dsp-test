# Copyright 2021 kubeflow.org
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from kfp import dsl
from kfp_tekton.compiler import TektonCompiler
from kfp import components


class Coder:
    def empty(self):
        return ""


TektonCompiler._get_unique_id_code = Coder.empty


@components.func_to_container_op
def add_numbers(a: int, b: int) -> int:
    print(a + b)
    return a + b


@components.func_to_container_op
def print_number(a: int) -> int:
    print(a)
    return a


@components.func_to_container_op
def notify_success():
    print('SUCCESS!')


@components.func_to_container_op
def notify_failure():
    print('FAILED!')


@components.func_to_container_op
def produce_numbers(n: int) -> list:
    import random
    rl = random.sample(range(0, 1000), n)
    print(rl)
    return rl


@components.func_to_container_op
def produce_number() -> int:
    import random
    rn = random.randrange(0, 1000)
    print(rn)
    return rn


@dsl.pipeline(name='conditions-and-loops')
def conditions_and_loops(n: int = 3, threshold: int = 20):
    produce_numbers_task = produce_numbers(n)
    with dsl.ParallelFor(produce_numbers_task.output) as loop_item:
        add_numbers_task = add_numbers(loop_item, 10)
        print_number_task = print_number(add_numbers_task.output)
        with dsl.Condition(print_number_task.output > threshold):
            notify_success()

        with dsl.Condition(print_number_task.output <= threshold):
            notify_failure()


if __name__ == '__main__':
    TektonCompiler().compile(conditions_and_loops, __file__.replace('.py', '.yaml'))
